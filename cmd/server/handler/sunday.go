package handler

import (
	"bytes"
	"devopegin/internal/domain"
	"devopegin/internal/sunday"
	"devopegin/pkg/web"
	"encoding/json"
	"errors"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ErrFileInvalid  = errors.New("file is invalid")
	ErrFileRequired = errors.New("file is required")
	ErrJsonRequired = errors.New("json data is required")
	ErrJsonData     = errors.New("json data is incorrect")
	ErrInternal     = errors.New("an internal error has occurred")
)

type Sunday struct {
	sundayService sunday.IService
}

func NewSunday(service sunday.IService) *Sunday {
	return &Sunday{
		sundayService: service,
	}
}

func (s *Sunday) GenerateDoc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Get file
		formFile, _, err := ctx.Request.FormFile("file")
		if err != nil {
			web.Error(ctx, http.StatusBadRequest, ErrFileRequired.Error())
			return
		}
		file, err := ioutil.ReadAll(formFile)
		if err != nil {
			web.Error(ctx, http.StatusInternalServerError, ErrInternal.Error())
			return
		}
		readerFile := bytes.NewReader(file)

		//Get Json Data
		jsonString := ctx.Request.Form.Get("sundayForm")
		if jsonString == "" {
			web.Error(ctx, http.StatusBadRequest, ErrFileRequired.Error())
			return
		}
		var sundayForm domain.SundayForm
		err = json.Unmarshal([]byte(jsonString), &sundayForm)
		if err != nil {
			web.Error(ctx, http.StatusBadRequest, ErrJsonData.Error())
			return
		}
		//Get image from path
		image, err := os.ReadFile("images/opegin.jpg")
		if err != nil {
			web.Error(ctx, http.StatusInternalServerError, ErrInternal.Error())
			return
		}

		buffer, err := s.sundayService.GenerateDocument(ctx, readerFile, domain.Sunday{
			Month:               sundayForm.Month,
			Year:                sundayForm.Year,
			Responsible:         sundayForm.Responsible,
			ImmediateBoss:       sundayForm.ImmediateBoss,
			EntryTime:           sundayForm.EntryTime,
			SundayEntryTime:     sundayForm.SundayEntryTime,
			SundayDepartureTime: sundayForm.SundayDepartureTime,
			Justification:       sundayForm.Justification,
			CompanyImage:        image,
		})
		if err != nil {
			switch {
			case errors.Is(err, sunday.ErrInvalidDocument):
				web.Error(ctx, http.StatusBadRequest, ErrFileInvalid.Error())
			default:
				web.Error(ctx, http.StatusInternalServerError, ErrInternal.Error())
			}
			return
		}
		downloadName := time.Now().UTC().Format("data-20060102150405.xlsx")
		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Disposition", "attachment; filename="+downloadName)
		//ctx.Data(http.StatusOK, http.DetectContentType(buffer.Bytes()), buffer.Bytes())
		ctx.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())
	}
}
