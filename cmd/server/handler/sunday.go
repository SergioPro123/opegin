package handler

import (
	"bytes"
	"devopegin/internal/sunday"
	"devopegin/pkg/web"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrFileRequired = errors.New("file is required")
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
		formFile, header, err := ctx.Request.FormFile("file")
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

		s.sundayService.GenerateDocument(ctx, readerFile)

		web.Success(ctx, http.StatusOK, fmt.Sprintf("'%s' uploaded!", header.Filename))
	}
}
