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
	if err := eCtx.Bind(user); err != nil {
		return err
	}
	user.Nickname = nickname
	_, err = DBS.ExecContext(eCtx.Request().Context(), `
		insert into 
		    actor (nickname,fullname,about,email) 
		values ($1,$2,$3,$4)
		`,
		user.Nickname, user.Fullname, user.About, user.Email)
	if err == nil {
		return eCtx.JSON(201, user)
	}
	rows, err := DBS.Query(`
		select nickname,fullname,about,email 
		from actor 
		where  nickname = $1 or email = $2
		`, user.Nickname, user.Email)

	defer rows.Close()
	var users []forms.User
	for rows.Next() {
		var rowUser forms.User
		if err = rows.Scan(
			&rowUser.Nickname,
			&rowUser.Fullname,
			&rowUser.About,
			&rowUser.Email,
		); err != nil {
			return eCtx.JSON(500, forms.Error{Message: "Cannot get user" + err.Error()})
		}
		users = append(users, rowUser)
	}

	return eCtx.JSON(409, users)
}

func GetUserProfile(eCtx echo.Context) error {
	nickname, err := getSlugUsername(eCtx)
	user := new(forms.User)
	if err = eCtx.Bind(user); err != nil {
		return err
	}
	user.Nickname = nickname
	err = DBS.QueryRowContext(eCtx.Request().Context(), ` 
			select fullname,about,email from actor where  nickname = $1`, user.Nickname).
		Scan(&user.Fullname, &user.About, &user.Email)

	if err != nil {
		return eCtx.JSON(404, forms.Error{Message: "Not found"})
	}

	return eCtx.JSON(200, user)
}

func UpdateUserProfile(eCtx echo.Context) error {
	nickname, err := getSlugUsername(eCtx)
	if err != nil {
		return err
	}
	user := new(forms.User)
	if err := eCtx.Bind(user); err != nil {
		return err
	}
	user.Nickname = nickname

	rows, err := DBS.QueryContext(eCtx.Request().Context(),
		`select nickname from actor 
                where nickname = $1
		`, user.Nickname, user.Email)
	if err != nil || rows == nil {
		return eCtx.JSON(404, forms.Error{Message: "none such user"})
	}
	defer rows.Close()
	if rows.Next() {
		return eCtx.JSON(404, forms.Error{Message: err.Error()})
	}
	_, err = DBS.ExecContext(eCtx.Request().Context(), `
		update actor 
		set fullname = $2,
		    about = $3,
		    email = $4 
		where nickname = $1
		`,
		user.Nickname, user.Fullname, user.About, user.Email)
	if err != nil {
		return eCtx.JSON(409, forms.Error{Message: "new params don't suit"})

	}
	return eCtx.JSON(200, user)
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
