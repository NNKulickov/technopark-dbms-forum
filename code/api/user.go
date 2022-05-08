package api

import (
	"github.com/NNKulickov/technopark-dbms-forum/forms"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateUser(eCtx echo.Context) error {
	nickname, err := getSlugUsername(eCtx)
	if err != nil {
		return err
	}
	user := new(forms.User)
	user.Nickname = nickname

	return eCtx.JSON(201, user)
}

func GetUserProfile(eCtx echo.Context) error {
	_, err := getSlugUsername(eCtx)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserProfile(eCtx echo.Context) error {
	_, err := getSlugUsername(eCtx)
	if err != nil {
		return err
	}
	return nil
}

func getSlugUsername(eCtx echo.Context) (string, error) {
	username := eCtx.Param(usernameSlug)
	if username == "" {
		eCtx.Logger().Error("cannot get username")
		return username, eCtx.JSON(
			http.StatusBadRequest,
			forms.Error{Message: "None user"})
	}
	return username, nil
}
