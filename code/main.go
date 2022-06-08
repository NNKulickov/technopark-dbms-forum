package main

import (
	"database/sql"
	"fmt"
	"github.com/NNKulickov/technopark-dbms-forum/api"
	_ "github.com/NNKulickov/technopark-dbms-forum/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"io/ioutil"
	"log"
	"os"
)

const initialScriptPath = "./DBScript/initial.sql"

func main() {
	e := echo.New()
	e.Debug = true
	e.GET("/docs/*", echoSwagger.WrapHandler)
	api.DBS = initDB(initialScriptPath)
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: `{"time":"${time_unix}",` +
				`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
				`"method":"${method}","uri":"${uri}",` +
				"\n",
		},
	))
	api.InitRoutes(e.Group("/api"))
	log.Fatal(e.Start("0.0.0.0:5000"))
}

func initDB(initDBPath string) *sql.DB {
	connectString := fmt.Sprintf(
		`host=%s
				port=5432
				user=%s
				password=%s
				dbname=%s
				sslmode=disable
				TimeZone=Europe/Moscow`,
		os.Getenv("DB_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	dbs, err := sql.Open("postgres", connectString)
	dbs.SetMaxOpenConns(5)
	dbs.SetMaxIdleConns(0)
	if err != nil {
		log.Fatal(err)
	}

	if err = dbs.Ping(); err != nil {
		log.Fatal(err)
	}
	sql, err := ioutil.ReadFile(initDBPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = dbs.Exec(string(sql))
	if err != nil {
		log.Fatal(err)
	}

	return dbs
}
