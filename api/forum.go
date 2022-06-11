package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/NNKulickov/technopark-dbms-forum/forms"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func CreateForum(eCtx echo.Context) error {
	forum := forms.PostForum{}
	if err := eCtx.Bind(&forum); err != nil {
		fmt.Println("CreateForum (1):", err)
		return err
	}
	ctx := eCtx.Request().Context()
	user, err := getUserFromDb(ctx, forum.User)
	if err != nil {
		return eCtx.JSON(404, forms.Error{
			Message: "none such user"})
	}
	_, err = DBS.ExecContext(ctx, `
		insert into forum (slug,title,host)
		values ($1,$2,$3)
	`,
		forum.Slug,
		forum.Title,
		user.Nickname,
	)
	if err == nil {
		return eCtx.JSON(201, forms.ForumResult{
			Title:   forum.Title,
			User:    user.Nickname,
			Slug:    forum.Slug,
			Posts:   0,
			Threads: 0,
		})
	}
	if res, err := getForum(ctx, forum.Slug); err == nil {
		return eCtx.JSON(http.StatusConflict, res)
	}
	fmt.Println("CreateForum (2):", err)
	return eCtx.JSON(404, forms.Error{
		Message: "none such user"})
}

func GetForumDetails(eCtx echo.Context) error {
	forumParam := eCtx.Param(forumSlug)
	ctx := eCtx.Request().Context()

	forum, err := getForum(ctx, forumParam)
	if err != nil {
		fmt.Println("CreateForum (1):", err)
		return eCtx.JSON(404, forms.Error{
			Message: "none such forum"})
	}

	return eCtx.JSON(200, forum)

}

func CreateForumThread(eCtx echo.Context) error {
	slug := eCtx.Param(forumSlug)
	thread := forms.ThreadForm{}
	ctx := eCtx.Request().Context()
	if err := eCtx.Bind(&thread); err != nil {
		fmt.Println("CreateForumThread (1):", err)
		return err
	}
	forum, err := getForum(ctx, slug)
	if err != nil {
		fmt.Println("CreateForumThread (2):", err)
		return eCtx.JSON(404, forms.Error{
			Message: "none such user or forum"})
	}
	threadModel := forms.ThreadModel{
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   forum.Slug,
		Message: thread.Message,
		Slug:    sql.NullString{String: thread.Slug, Valid: true},
	}

	// try insert
	builder := strings.Builder{}
	builder.WriteString("insert into thread (title,author,forum,message,slug")
	if thread.Created != "" {
		builder.WriteString(",created")
	}
	builder.WriteString(") values ($1,$2,$3,$4,nullif($5,'')")
	if thread.Created != "" {

		builder.WriteString(fmt.Sprintf(",'%s'", thread.Created))
	}
	builder.WriteString(") returning id,created")
	if err = DBS.QueryRowContext(ctx, builder.String(),
		threadModel.Title,
		threadModel.Author,
		threadModel.Forum,
		threadModel.Message,
		threadModel.Slug,
	).
		Scan(
			&threadModel.Id,
			&threadModel.Created,
		); err == nil {
		return eCtx.JSON(201, forms.ThreadForm{
			Id:      threadModel.Id,
			Title:   threadModel.Title,
			Author:  threadModel.Author,
			Forum:   threadModel.Forum,
			Message: threadModel.Message,
			Slug:    threadModel.Slug.String,
			Votes:   threadModel.Votes,
			Created: strfmt.DateTime(threadModel.Created.UTC()).String(),
		})
	}
	// select if exists
	if err = DBS.QueryRowContext(eCtx.Request().Context(), `
		select
		id,
		title,
		author,
		forum,
		message,
		slug,
		created
		from thread 
		where lower(slug) = lower($1)`,
		threadModel.Slug.String,
	).
		Scan(
			&threadModel.Id,
			&threadModel.Title,
			&threadModel.Author,
			&threadModel.Forum,
			&threadModel.Message,
			&threadModel.Slug,
			&threadModel.Created,
		); err == nil {
		return eCtx.JSON(409, forms.ThreadForm{
			Id:      threadModel.Id,
			Title:   threadModel.Title,
			Author:  threadModel.Author,
			Forum:   threadModel.Forum,
			Message: threadModel.Message,
			Slug:    threadModel.Slug.String,
			Created: strfmt.DateTime(threadModel.Created.UTC()).String(),
		})
	}
	fmt.Println("CreateForumThread not found (3)", err)
	return eCtx.JSON(404, forms.Error{
		Message: "none such user or forum"})
}

func GetForumUsers(eCtx echo.Context) error {
	slug := eCtx.Param(forumSlug)
	ctx := eCtx.Request().Context()
	forum, err := getForum(ctx, slug)
	if err != nil {
		fmt.Println("GetForumUsers (0):", err)
		return eCtx.JSON(http.StatusNotFound, forms.Error{
			Message: "none such forum"})
	}
	slug = forum.Slug
	limit := 0
	desc := false
	limitString := eCtx.QueryParam("limit")
	descString := eCtx.QueryParam("desc")
	if limit, err = strconv.Atoi(limitString); err != nil {
		fmt.Println("GetForumUsers (1):", err)
	}
	if desc, err = strconv.ParseBool(descString); err != nil {
		fmt.Println("GetForumUsers (2):", err)
	}
	users := forms.UserFilter{
		Limit: limit,
		Desc:  desc,
		Since: eCtx.QueryParam("since"),
	}
	usersQuery := `
		select a.nickname,a.fullname,a.about,a.email from post p
			join actor a on lower(p.author) = lower(a.nickname) 
		where lower(p.forum) = lower($1)
		union
		select a.nickname,a.fullname,a.about,a.email from thread t
			join actor a on lower(t.author) = lower(a.nickname)
		where lower(t.forum) = lower($1)`

	build := strings.Builder{}
	build.WriteString(fmt.Sprintf(`
			select nickname,fullname,about,email from ( %s ) users
	`, usersQuery))
	if users.Since != "" {
		if users.Desc {
			build.WriteString(fmt.Sprintf(` where lower(nickname) collate "C" <  lower('%s') collate "C"`, users.Since))

		} else {
			build.WriteString(fmt.Sprintf(` where lower(nickname) collate "C" >  lower('%s') collate "C"`, users.Since))

		}
	}
	build.WriteString(` order by lower(nickname) collate "C"`)
	if users.Desc {
		build.WriteString(" desc")
	}
	build.WriteString(" limit nullif($2,0)")
	if rows, err := DBS.QueryContext(ctx, build.String(), slug, users.Limit); err == nil {
		usersResponse := make([]forms.User, 0, 100)
		for rows.Next() {
			user := forms.User{}
			if err = rows.
				Scan(
					&user.Nickname,
					&user.Fullname,
					&user.About,
					&user.Email,
				); err != nil {
				fmt.Println("GetForumUsers (2):", err)

				return eCtx.JSON(http.StatusInternalServerError, forms.Error{
					Message: "smth wrong"})
			}
			usersResponse = append(usersResponse, user)
		}
		return eCtx.JSON(http.StatusOK, usersResponse)
	}
	return eCtx.JSON(http.StatusNotFound, forms.Error{
		Message: "none such forum"})
}

func GetForumThreads(eCtx echo.Context) error {
	slug := eCtx.Param(forumSlug)
	ctx := eCtx.Request().Context()
	forum := ""
	if err := DBS.QueryRowContext(ctx,
		`select slug from forum where lower(slug) = lower($1)`,
		slug).Scan(&forum); err != nil {
		return eCtx.JSON(http.StatusNotFound, forms.Error{
			Message: "none forum"})
	}
	slug = forum
	limit, err := strconv.Atoi(eCtx.QueryParam("limit"))
	if err != nil {
		fmt.Println("err:", err)
		limit = 0
	}
	desc, err := strconv.ParseBool(eCtx.QueryParam("desc"))
	if err != nil {
		fmt.Println("err:", err)
		desc = false
	}
	threadsFilter := forms.ThreadFilter{
		Limit: limit,
		Since: eCtx.QueryParam("since"),
		Desc:  desc,
	}
	if err := eCtx.Bind(&threadsFilter); err != nil {
		fmt.Println("GetForumUsers (3):", err)
		return err
	}
	build := strings.Builder{}
	build.WriteString(`
		select id, title, author, forum, message,
			slug, created,votes from thread
			where lower(forum) = lower($1)`)
	if threadsFilter.Since != "" {
		if desc {
			build.WriteString(fmt.
				Sprintf(
					` and created <= '%s'`,
					threadsFilter.Since,
				),
			)
		} else {
			build.WriteString(fmt.
				Sprintf(
					` and created >= '%s'`,
					threadsFilter.Since,
				),
			)
		}
	}
	build.WriteString(" order by created")
	if threadsFilter.Desc {
		build.WriteString(" desc")
	}
	build.WriteString(" limit nullif($2,0)")
	if rows, err := DBS.
		QueryContext(
			ctx,
			build.String(),
			slug,
			threadsFilter.Limit); err == nil {
		threadsResponse := make([]forms.ThreadForm, 0, 100)
		for rows.Next() {
			thread := forms.ThreadModel{}
			if err = rows.
				Scan(
					&thread.Id,
					&thread.Title,
					&thread.Author,
					&thread.Forum,
					&thread.Message,
					&thread.Slug,
					&thread.Created,
					&thread.Votes,
				); err != nil {
				fmt.Println("GetForumUsers (4):", err)

				return eCtx.JSON(http.StatusInternalServerError, forms.Error{
					Message: "smth wrong"})
			}
			threadsResponse = append(threadsResponse, forms.ThreadForm{
				Id:      thread.Id,
				Title:   thread.Title,
				Author:  thread.Author,
				Forum:   thread.Forum,
				Message: thread.Message,
				Slug:    thread.Slug.String,
				Votes:   thread.Votes,
				Created: strfmt.DateTime(thread.Created.UTC()).String(),
			})
		}
		return eCtx.JSON(http.StatusOK, threadsResponse)
	}
	return eCtx.JSON(http.StatusNotFound, forms.Error{
		Message: "none such forum"})
}

func getForum(ctx context.Context, slug string) (forms.ForumResult, error) {
	forum := forms.ForumResult{}

	err := DBS.QueryRowContext(ctx, `
	with 
	    threads as (
			select count(*) threads from thread where lower(forum) = lower($1)
		),
	    posts as (
	    	select count(*) posts from post where lower(forum)= lower($1)
	    )
	select f.slug,f.title,f.host,t.threads,p.posts from threads t,posts p, forum f where lower(f.slug) = lower($1);
	`, slug).
		Scan(
			&forum.Slug,
			&forum.Title,
			&forum.User,
			&forum.Threads,
			&forum.Posts,
		)
	return forum, err
}
