package api

import (
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
)

const initialScriptPath = "./DBScript/initial.sql"

func GetServiceStatus(eCtx echo.Context) error { return nil }
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
