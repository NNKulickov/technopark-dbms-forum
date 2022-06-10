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
	res := forms.ForumResult{}
	if err = DBS.QueryRowContext(ctx, `
		select f.slug,f.title,f.host,count(th) threads, count(p) posts
		from forum f
			left join thread th 
			    on th.forum = f.slug
			left join post p 
			    on p.forum = f.slug and p.threadid = th.id
		where lower(f.slug) = lower($1) group by f.slug
	`, forum.Slug).
		Scan(
			&res.Slug,
			&res.Title,
			&res.User,
			&res.Threads,
			&res.Posts,
		); err == nil {
		return eCtx.JSON(409, res)
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
		fmt.Println("CreateForum (2):", err)
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
	users := forms.UserFilter{}
	if err := eCtx.Bind(&users); err != nil {
		fmt.Println("GetForumUsers (1):", err)
		return err
	}
	build := strings.Builder{}
	build.WriteString("with users(nickname,fullname,about,email) as (")
	addSourceUser(&build, "post")
	addSinceUser(&build, users.Since)
	addSourceUser(&build, "thread")
	addSinceUser(&build, users.Since)
	build.WriteString(`
		select nickname,fullname,about,email from users 
		    ORDER BY lower(nickname) 
		    limit nullif($2,0)
		`)
	if users.Desc {
		build.WriteString(" Desc")
	}
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

func addSourceUser(builder *strings.Builder, src string) {

	builder.WriteString(fmt.Sprintf(`
		 	select a.nickname, a.fullname, a.about, a.email from actor a
				join %s src on a.nickname = src.author
				where src.forum = $1`, src))
}

func addSinceUser(builder *strings.Builder, since string) {
	if since != "" {
		builder.WriteString(fmt.Sprintf(` and src.nickname > '%s'`, since))
	}
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
		select id,title,author,forum, message,
			slug, created from thread
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
	build.WriteString(" ORDER BY created")
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
		select f.slug,f.title,f.host,count(th) threads, count(p) posts
		from forum f
			left join thread th 
			    on th.forum = f.slug
			left join post p 
			    on p.forum = f.slug and p.threadid = th.id
		where lower(f.slug) = lower($1) group by f.slug
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
