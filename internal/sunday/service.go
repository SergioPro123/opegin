package sunday

import (
	"bytes"
	"devopegin/internal/domain"
	"devopegin/pkg/utils"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var (
	ErrInvalidDocument = errors.New("the document does not meet the conditions to be processed")
	ErrInternal        = errors.New("an internal error has occurred")
)

var (
	BorderDefauld = []excelize.Border{
		{
			Type:  "left",
			Color: "000000",
			Style: 1,
		},
		{
			Type:  "top",
			Color: "000000",
			Style: 1,
		},
		{
			Type:  "right",
			Color: "000000",
			Style: 1,
		},
		{
			Type:  "bottom",
			Color: "000000",
			Style: 1,
		},
	}
	FontWithBold = &excelize.Font{
		Bold: true,
		Size: 10,
	}
	FontWithoutBold = &excelize.Font{
		Bold: false,
		Size: 10,
	}
	FontWithBoldAndWhite = &excelize.Font{
		Bold:  true,
		Size:  10,
		Color: "ffffff",
	}
	fillColorBlue = excelize.Fill{
		Color:   []string{"#000066"},
		Type:    "pattern",
		Pattern: 1,
	}

	AlignmentDefault = &excelize.Alignment{
		Horizontal: "center",
		Vertical:   "center",
		WrapText:   true,
	}
	AlignmentRight = &excelize.Alignment{
		Horizontal: "right",
		Vertical:   "center",
		WrapText:   true,
	}
	AlignmentLeft = &excelize.Alignment{
		Horizontal: "left",
		Vertical:   "center",
		WrapText:   true,
	}
	AlignmentLeft_Top = &excelize.Alignment{
		Horizontal: "left",
		Vertical:   "top",
		WrapText:   true,
	}
)

type IService interface {
	GenerateDocument(ctx *gin.Context, sundays io.Reader, sundayForm domain.Sunday) (*bytes.Buffer, error)
}

type service struct {
	repository             IRepository
	StartTimeForExtraHours time.Time
}

func NewService(repository IRepository) IService {
	time, _ := time.Parse("15:04", "21:00")
	return &service{
		repository: repository,
		//In Colombia, overtime is paid from 9 P.M.
		StartTimeForExtraHours: time,
	}
}

func (s *service) GenerateDocument(ctx *gin.Context, sundays io.Reader, sundayForm domain.Sunday) (*bytes.Buffer, error) {
	employees, err := s.LoadEmployeesFromExcel(sundays)
	if err != nil {
		return nil, err
	}
	buffer, err := s.EmployeesToExcelCalculed(employees, sundayForm)
	if err != nil {
		return nil, ErrInternal
	}
	return buffer, nil
}

func (s *service) LoadEmployeesFromExcel(sundays io.Reader) (employees []*domain.Employee, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = ErrInternal
		}
	}()
	doc, err := excelize.OpenReader(sundays)
	if err != nil {
		fmt.Println(err)
		return []*domain.Employee{}, ErrInvalidDocument
	}
	defer func() {
		// Close the spreadsheet.
		if err := doc.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	//Get name from first sheet
	firstSheetName := doc.WorkBook.Sheets.Sheet[0].Name

	rows, err := doc.GetRows(firstSheetName)
	if err != nil {
		fmt.Println(err)
		return []*domain.Employee{}, ErrInternal
	}
	//Iterates over all dates y save in repository
	rowDates := rows[1][4:]
	countDates := 1
	for i := 0; i < len(rowDates); i++ {
		//the date is parsed
		date, err := time.Parse("01-02-06", rowDates[i])
		if err != nil {
			return []*domain.Employee{}, ErrInvalidDocument
		}
		//the time is parsed
		hour, err := strconv.Atoi(rowDates[i+1])
		if err != nil {
			return []*domain.Employee{}, ErrInvalidDocument
		}
		//Save in repository
		s.repository.AddExtraHour(domain.ExtraHour{
			ID:            countDates,
			Date:          date,
			NumberOfHours: hour,
		})
		i++ //x2 (this + for)
		countDates++
	}

	//Each employee iterates
	rowsEmployees := rows[2:]
	for _, rowEmployee := range rowsEmployees {
		var employee *domain.Employee = &domain.Employee{}
		employees = append(employees, employee)

		//Get document from employee
		employee.Document = rowEmployee[0]
		//Get name from employee
		employee.Name = rowEmployee[1]
		//check if locality exists, if not then insert
		namePosition := rowEmployee[2]
		//Get group from employee
		group := 0
		if len(rowEmployee) >= 4 {
			if rowEmployee[3] == "" {
				rowEmployee[3] = "0"
			}
			group, err = strconv.Atoi(rowEmployee[3])
			if err != nil {
				return []*domain.Employee{}, ErrInvalidDocument
			}
		}

		employee.Group = group
		location := s.repository.GetPosition(namePosition)
		if location == nil {
			s.repository.AddPosition(domain.Position{
				Name: namePosition,
			})
			location = s.repository.GetPosition(namePosition)
		}
		//Add location
		employee.Position = location
		//iterates over the overtime applied
		employee.ExtraHours = []*domain.ExtraHour{}
		if !(len(rowEmployee) >= 5) {
			continue
		}
		rowApplied := rowEmployee[4:]
		countApplied := 1
		for i := 0; i < len(rowApplied); i++ {
			//If it is different from empty, then it is because overtime applies
			if rowApplied[i] != "" {
				extraHour := s.repository.GetExtraHour(countApplied)
				if extraHour == nil {
					return []*domain.Employee{}, ErrInternal
				}
				employee.ExtraHours = append(employee.ExtraHours, extraHour)
			}
			//the next item is ignored due to the merged cell
			i++ //x2 (this + for)
			countApplied++
		}
	}
	return employees, nil
}

func (s *service) EmployeesToExcelCalculed(employees []*domain.Employee, sundayForm domain.Sunday) (*bytes.Buffer, error) {
	doc := excelize.NewFile()

	//Delete default sheet
	err := doc.DeleteSheet(doc.WorkBook.Sheets.Sheet[0].Name)
	if err != nil {
		return nil, err
	}
	//Get difference hour between hour extra and hour pay legal
	extraHourDifferenceToPay := GetDifferenceHours(s.StartTimeForExtraHours, sundayForm.EntryTime.Time)

	for _, employee := range employees {
		nameSheet := utils.GetFirstNameAndFirstLastName(employee.Name)

		indexDoc, err := doc.NewSheet(nameSheet)
		if err != nil {
			return nil, ErrInternal
		}
		doc.SetActiveSheet(indexDoc)

		//Set the width of the columns
		var widthOdCol map[string]float64 = map[string]float64{
			"A": 0.58, "B": 10.71, "C": 14.71, "D": 10.71, "E": 10.71, "F": 10.71, "G": 8.14, "H": 8.14, "I": 10.71,
			"J": 10.71, "K": 10.71, "L": 10.71, "M": 10.71, "N": 10.71, "O": 10.71, "P": 0.58}
		for key, value := range widthOdCol {
			err = doc.SetColWidth(nameSheet, key, key, value)
			if err != nil {
				return nil, ErrInternal
			}
		}

		//Merge cells
		var mergeCols map[string]string = map[string]string{
			"B1": "O1",
			"B2": "C4", "D2": "M4", "N2": "O4",
			"B5": "O5",
			"B6": "D7", "E6": "F7", "G6": "K7", "L6": "O7",
			"B8": "G9", "H8": "K9", "L8": "O9",
			"B10": "G12", "H10": "K12", "L10": "O12",
			"B13": "O13",
			"I14": "L14", "M14": "O14",
			"B25": "O25",
			"B26": "O27",
			"B28": "H28", "I28": "O28",
		}
		//Cell combinations are added according to the overtime you have (minimum 10 records in the table)
		totalColOfExtraHours := len(employee.ExtraHours)
		if totalColOfExtraHours <= 10 {
			totalColOfExtraHours = 10
		}
		//Add combinations of table extra hours
		for i := 15; i < (15 + totalColOfExtraHours); i++ {
			mergeCols["I"+strconv.Itoa(i)] = "L" + strconv.Itoa(i)
			mergeCols["M"+strconv.Itoa(i)] = "O" + strconv.Itoa(i)
		}
		for key, value := range mergeCols {
			err = doc.MergeCell(nameSheet, key, value)
			if err != nil {
				return nil, err
			}
		}

		//Set the height of the rows
		var heightRow map[int]float64 = map[int]float64{
			12:                              5,
			14:                              30.75,
			(1 + 15 + totalColOfExtraHours): 20.25,
			(2 + 15 + totalColOfExtraHours): 20.25,
			(3 + 15 + totalColOfExtraHours): 23.25,
			(4 + 15 + totalColOfExtraHours): 9,
		}
		//Add height of row extra hours table
		for i := 15; i < 15+totalColOfExtraHours; i++ {
			heightRow[i] = 37.5
		}
		for key, value := range heightRow {
			err = doc.SetRowHeight(nameSheet, key, value)
			if err != nil {
				return nil, ErrInternal
			}
		}
		// *************************** Add Styles ***************************
		//Add style to Title
		style, err := doc.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Bold: true,
				Size: 16},
			Border:    BorderDefauld,
			Alignment: AlignmentDefault,
		})
		err = doc.SetCellStyle(nameSheet, "B2", "M4", style)
		if err != nil {
			return nil, ErrInternal
		}
		style, err = doc.NewStyle(&excelize.Style{
			Font:      FontWithoutBold,
			Border:    BorderDefauld,
			Alignment: AlignmentRight,
		})
		err = doc.SetCellStyle(nameSheet, "N2", "O4", style)
		if err != nil {
			return nil, ErrInternal
		}
		//Add style to Header
		style, err = doc.NewStyle(&excelize.Style{
			Font:      FontWithBold,
			Border:    BorderDefauld,
			Alignment: AlignmentLeft,
		})

		err = doc.SetCellStyle(nameSheet, "B6", "O12", style)
		if err != nil {
			return nil, ErrInternal
		}
		//Add style to table header
		style, err = doc.NewStyle(&excelize.Style{
			Font:      FontWithBoldAndWhite,
			Border:    BorderDefauld,
			Fill:      fillColorBlue,
			Alignment: AlignmentDefault,
		})
		err = doc.SetCellStyle(nameSheet, "B14", "O14", style)
		if err != nil {
			return nil, ErrInternal
		}
		//Add style to table content
		style, err = doc.NewStyle(&excelize.Style{
			Border:    BorderDefauld,
			Alignment: AlignmentDefault,
		})

		err = doc.SetCellStyle(nameSheet, "B15", "O"+strconv.Itoa(14+totalColOfExtraHours), style)
		if err != nil {
			return nil, ErrInternal
		}
		//Add style to table footer
		style, err = doc.NewStyle(&excelize.Style{
			Border:    BorderDefauld,
			Alignment: AlignmentLeft,
			Font:      FontWithBold,
		})

		err = doc.SetCellStyle(nameSheet, "B"+strconv.Itoa(14+totalColOfExtraHours), "O"+strconv.Itoa(4+14+totalColOfExtraHours), style)
		if err != nil {
			return nil, ErrInternal
		}

		style, err = doc.NewStyle(&excelize.Style{
			Border:    BorderDefauld,
			Alignment: AlignmentLeft_Top,
			Font:      FontWithBold,
		})

		err = doc.SetCellStyle(nameSheet, "B"+strconv.Itoa(15+totalColOfExtraHours), "O"+strconv.Itoa(3+15+totalColOfExtraHours), style)
		if err != nil {
			return nil, ErrInternal
		}
		// ******************************* Add images to cells *******************************
		err = doc.AddPictureFromBytes(nameSheet, "B2", "opegin", ".jpg", sundayForm.CompanyImage, &excelize.GraphicOptions{
			AutoFit: true,
			ScaleX:  1.7,
			ScaleY:  0.9,
		})
		if err != nil {
			return nil, ErrInternal
		}
		// ******************************* Add values to cells *******************************
		valuesCell := map[string]string{
			//*Title
			"D2": "Novedades de Trabajo Nocturno, Horas Extra, Trabajo Dominical y Festivo",
			"N2": "Código: F-GA-06 \n Fecha: 08/02/17 \n Versión: 2",
			//*Header
			"B6":  "Mes: " + utils.CapitalizeWords(sundayForm.Month),
			"E6":  "Año: " + utils.CapitalizeWords(sundayForm.Year),
			"G6":  "Novedad Realizada por: " + utils.CapitalizeWords(sundayForm.Responsible.Name),
			"L6":  "Cargo: " + utils.CapitalizeWords(sundayForm.Responsible.Position.Name),
			"B8":  "Nombre del Trabajador: " + utils.CapitalizeWords(employee.Name),
			"H8":  "No Identificación: " + utils.CapitalizeWords(employee.Document),
			"L8":  "Cargo: " + utils.CapitalizeWords(employee.Position.Name),
			"B10": "Nombre de Jefe Inmediato: " + utils.CapitalizeWords(sundayForm.ImmediateBoss.Name),
			"H10": "Localidad: " + utils.CapitalizeWords(sundayForm.ImmediateBoss.Location),
			"L10": "Departamento: " + utils.CapitalizeWords(sundayForm.ImmediateBoss.Department),
			//*Table Extra Hours (Header)
			"B14": "Fecha",
			"C14": "Hora de Entrada",
			"D14": "Hora de Salida",
			"E14": "Total Horas Diurnas",
			"F14": "Total Horas Nocturnas",
			"G14": "Domingo",
			"H14": "Festivo",
			"I14": "Justificación",
			"M14": "Firma del Trabajador",
			//*Table footer
			"B" + strconv.Itoa((15 + totalColOfExtraHours)):     "*** Hora Diurna:  de 6:00 am a 10:00 pm",
			"B" + strconv.Itoa((1 + 15 + totalColOfExtraHours)): "Observaciones:",
			"B" + strconv.Itoa((3 + 15 + totalColOfExtraHours)): "VoBo Ingeniero  Residente - Coordinador Departamental:",
			"I" + strconv.Itoa((3 + 15 + totalColOfExtraHours)): "VoBo Director Operativo:",
		}
		//Calculations are added according to your applied overtime
		numberCellTable := 15
		for _, extraHour := range employee.ExtraHours {
			CellRowNumber := strconv.Itoa(numberCellTable)
			valuesCell["B"+CellRowNumber] = extraHour.Date.Format("2006-01-02")
			//is it sunday?
			if extraHour.Date.Weekday() == 0 {
				valuesCell["C"+CellRowNumber] = sundayForm.SundayEntryTime.Format("03:04 PM")
				valuesCell["D"+CellRowNumber] = sundayForm.SundayDepartureTime.Format("03:04 PM")
				valuesCell["G"+CellRowNumber] = fmt.Sprintf("%.2f",
					getHourFloat(sundayForm.SundayDepartureTime.Time)-getHourFloat(sundayForm.SundayEntryTime.Time))

			} else {
				valuesCell["C"+CellRowNumber] = sundayForm.EntryTime.Format("03:04 PM")
				valuesCell["D"+CellRowNumber] = sundayForm.EntryTime.Add(time.Hour * time.Duration(extraHour.NumberOfHours)).Format("03:04 PM")

				diferenceExtraHourMade := extraHourDifferenceToPay - float64(extraHour.NumberOfHours)
				if diferenceExtraHourMade >= 0 {
					valuesCell["E"+CellRowNumber] = strconv.Itoa(extraHour.NumberOfHours)
				} else {
					valuesCell["E"+CellRowNumber] = fmt.Sprintf("%.2f", extraHourDifferenceToPay)
					valuesCell["F"+CellRowNumber] = fmt.Sprintf("%.2f", math.Abs(diferenceExtraHourMade))
				}
			}

			valuesCell["I"+CellRowNumber] = sundayForm.Justification

			numberCellTable++
		}
		//Apply values to cell
		for key, value := range valuesCell {
			err = doc.SetSheetRow(nameSheet, key, &[]interface{}{value})
			if err != nil {
				return nil, ErrInternal
			}
		}

	}

	return doc.WriteToBuffer()
}

func GetDifferenceHours(firstHour, secondHour time.Time) float64 {
	return getHourFloat(firstHour) - getHourFloat(secondHour)
}

func getHourFloat(time time.Time) float64 {
	firstHourInt := time.Hour()
	firstMinuteInt := time.Minute()
	return float64(firstHourInt) + (float64(firstMinuteInt) / 60.0)
}
