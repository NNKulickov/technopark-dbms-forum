package main

import (
	"github.com/NNKulickov/technopark-dbms-forum/api"
	_ "github.com/NNKulickov/technopark-dbms-forum/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
)

func main() {
	e := echo.New()
	e.Debug = true
	e.GET("/docs/*", echoSwagger.WrapHandler)
	api.InitRoutes(e.Group("/api"))
	log.Fatal(e.Start("0.0.0.0:5000"))
}
