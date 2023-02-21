package handler

import (
	"devopegin/internal/sunday"
	"devopegin/pkg/web"
	"errors"
	"fmt"
	"log"
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

func (e *Sunday) GenerateDoc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			web.Error(ctx, http.StatusBadRequest, ErrFileRequired.Error())
			return
		}
		log.Println(file.Filename)

		err = ctx.SaveUploadedFile(file, file.Filename)
		if err != nil {
			fmt.Println(err.Error())
			web.Error(ctx, http.StatusInternalServerError, ErrInternal.Error())
			return
		}
		web.Success(ctx, http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	}
}
