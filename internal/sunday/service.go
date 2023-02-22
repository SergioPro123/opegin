package sunday

import (
	"devopegin/internal/domain"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var (
	ErrInvalidDocument = errors.New("the document does not meet the conditions to be processed")
	ErrInternal        = errors.New("an internal error has occurred")
)

type IService interface {
	GenerateDocument(ctx *gin.Context, sundays io.Reader) error
}

type service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	return &service{
		repository: repository,
	}
}

func (s *service) GenerateDocument(ctx *gin.Context, sundays io.Reader) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = ErrInternal
		}
	}()
	doc, err := excelize.OpenReader(sundays)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func() {
		// Close the spreadsheet.
		if err := doc.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var employees []*domain.Employee = []*domain.Employee{}
	//Get name from first sheet
	firstSheetName := doc.WorkBook.Sheets.Sheet[0].Name

	rows, err := doc.GetRows(firstSheetName)
	if err != nil {
		fmt.Println(err)
		return ErrInternal
	}
	//Iterates over all dates y save in repository
	rowDates := rows[1][4:]
	countDates := 1
	for i := 0; i < len(rowDates); i++ {
		//the date is parsed
		date, err := time.Parse("01-02-06", rowDates[i])
		if err != nil {
			return ErrInvalidDocument
		}
		//the time is parsed
		hour, err := strconv.Atoi(rowDates[i+1])
		if err != nil {
			return ErrInvalidDocument
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
		nameLocation := rowEmployee[2]
		//Get group from employee
		group := 0
		if len(rowEmployee) >= 4 {
			if rowEmployee[3] == "" {
				rowEmployee[3] = "0"
			}
			group, err = strconv.Atoi(rowEmployee[3])
			if err != nil {
				return ErrInvalidDocument
			}
		}

		employee.Group = group
		location := s.repository.GetLocation(nameLocation)
		if location == nil {
			s.repository.AddLocation(domain.Location{
				Name: nameLocation,
			})
			location = s.repository.GetLocation(nameLocation)
		}
		//Add location
		employee.Location = location
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
					return ErrInternal
				}
				employee.ExtraHours = append(employee.ExtraHours, extraHour)
			}
			//the next item is ignored due to the merged cell
			i++ //x2 (this + for)
			countApplied++
		}
	}
	return nil
}
