package routes

import (
	"devopegin/cmd/server/handler"
	"devopegin/internal/sunday"

	"github.com/gin-gonic/gin"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng *gin.Engine
	rg  *gin.RouterGroup
}

func NewRouter(eng *gin.Engine) Router {
	return &router{eng: eng}
}

func (r *router) MapRoutes() {
	r.setGroup()

	r.buildGenerateSunday()

}

func (r *router) setGroup() {
	r.rg = r.eng.Group("/api/v1")
}

func (r *router) buildGenerateSunday() {
	// Example
	repo := sunday.NewRepository()
	service := sunday.NewService(repo)
	handler := handler.NewSunday(service)
	r.rg.POST("/sundays", handler.GenerateDoc())
	r.eng.GET("/sundays", handler.GetView())
}
