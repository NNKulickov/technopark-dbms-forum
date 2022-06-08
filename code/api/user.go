package api

import (
	"context"
	"database/sql"
	"github.com/NNKulickov/technopark-dbms-forum/forms"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateUser(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
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
	rows, err := DBS.QueryContext(ctx, `
		select nickname,fullname,about,email 
		from actor 
		where  lower(nickname) = lower($1) or lower(email) = lower($2)
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
	user, err := getUserFromDb(eCtx.Request().Context(), nickname)
	if err != nil {
		return eCtx.JSON(404, forms.Error{Message: "Not found"})
	}

	return eCtx.JSON(200, user)
}

func getUserFromDb(ctx context.Context, nickname string) (forms.User, error) {
	user := forms.User{}
	err := DBS.QueryRowContext(ctx, ` 
			select nickname,fullname,about,email from actor where  lower(nickname) = lower($1)`, nickname).
		Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)

	return user, err
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
	userModel := forms.User{}

	if err = DBS.QueryRowContext(eCtx.Request().Context(),
		`select nickname,fullname,about,email from actor 
                where lower(nickname) = lower($1)
		`, user.Nickname).Scan(
		&userModel.Nickname,
		&userModel.Fullname,
		&userModel.About,
		&userModel.Email,
	); err == sql.ErrNoRows {
		return eCtx.JSON(404, forms.Error{Message: "none such user"})

	}
	if user.About == "" {
		user.About = userModel.About
	}
	if user.Email == "" {
		user.Email = userModel.Email
	}
	if user.Fullname == "" {
		user.Fullname = userModel.Fullname
	}
	if err = DBS.QueryRowContext(eCtx.Request().Context(), `
		update actor 
		set fullname = $2,
		    about = $3,
		    email = $4
		where lower(nickname) = lower($1)
		returning nickname,fullname,about,email
		`,
		user.Nickname, user.Fullname, user.About, user.Email).
		Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		); err != nil {
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
