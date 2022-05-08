package main

import (
	"database/sql"
	"fmt"
	"github.com/NNKulickov/technopark-dbms-forum/api"
	_ "github.com/NNKulickov/technopark-dbms-forum/docs"
	"github.com/labstack/echo/v4"
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
	api.InitRoutes(e.Group("/api"))
	log.Fatal(e.Start("0.0.0.0:5000"))
}

func initDB(initDBPath string) *sql.DB {
	connectString := fmt.Sprintf(
		`host=192.168.99.100
				port=5432
				user=%s
				password=%s
				dbname=%s
				sslmode=disable
				TimeZone=Europe/Moscow`,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	fmt.Println("hui")
	dbs, err := sql.Open("postgres", connectString)

	if err != nil {
		fmt.Println("hui1")
		log.Fatal(err)
	}

	if err = dbs.Ping(); err != nil {
		fmt.Println("hui2")
		log.Fatal(err)
	}

	sql, err := ioutil.ReadFile(initDBPath)
	if err != nil {
		fmt.Println("hui3")
		log.Fatal(err)
	}

	_, err = dbs.Exec(string(sql))
	if err != nil {
		fmt.Println("hui4")
		log.Fatal(err)
	}

	return dbs
}
