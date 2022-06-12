package api

import (
	"fmt"
	"github.com/NNKulickov/technopark-dbms-forum/forms"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
)

const initialScriptPath = "./db/db.sql"

func GetServiceStatus(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	status := forms.Status{}
	if err := DBS.QueryRowContext(ctx, `select count(*) from forum`).
		Scan(&status.Forum); err != nil {
		fmt.Println("GetServiceStatus (1) :", err)
		return err
	}

	if err := DBS.QueryRowContext(ctx, `select count(*) from post`).
		Scan(&status.Post); err != nil {
		fmt.Println("GetServiceStatus (2) :", err)
		return err
	}

	if err := DBS.QueryRowContext(ctx, `select count(*) from thread`).
		Scan(&status.Thread); err != nil {
		fmt.Println("GetServiceStatus (3) :", err)
		return err
	}

	if err := DBS.QueryRowContext(ctx, `select count(*) from actor`).
		Scan(&status.User); err != nil {
		fmt.Println("GetServiceStatus (4) :", err)
		return err
	}
	return eCtx.JSON(http.StatusOK, status)
}

func ClearServiceData(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	_, err := DBS.ExecContext(ctx, `
		truncate actor,forum,post,thread,vote
		`)
	sql, err := ioutil.ReadFile(initialScriptPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DBS.ExecContext(ctx, string(sql))
	if err != nil {
		log.Fatal(err)
	}
	return eCtx.JSON(http.StatusOK, "ok")
}
