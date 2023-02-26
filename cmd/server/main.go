package main

import (
	"devopegin/cmd/server/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	eng := gin.Default()

	eng.Static("/assets", "./public/assets")
	eng.LoadHTMLGlob("views/*.hbs")

	router := routes.NewRouter(eng)
	router.MapRoutes()

	if err := eng.Run(); err != nil {
		panic(err)
	}
}
