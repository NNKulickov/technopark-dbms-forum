package api

import (
	"database/sql"
	"fmt"
	"github.com/NNKulickov/technopark-dbms-forum/forms"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func CreateForum(eCtx echo.Context) error {
	forum := forms.PostForum{}
	if err := eCtx.Bind(&forum); err != nil {
		fmt.Println("CreateForum (1):", err)
		return err
	}
	ctx := eCtx.Request().Context()

	_, err := DBS.ExecContext(ctx, `
		insert into forum (slug,title,host) 
		values ($1,$2,$3)
	`,
		forum.Slug,
		forum.Title,
		forum.User,
	)
	if err == nil {
		return eCtx.JSON(201, forms.ForumResult{
			Title:   forum.Title,
			User:    forum.User,
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
		where f.slug = $1 group by f.slug
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
	forum := forms.ForumResult{}
	ctx := eCtx.Request().Context()

	err := DBS.QueryRowContext(ctx, `
		select f.slug,f.title,f.host,count(th) threads, count(p) posts
		from forum f
			left join thread th 
			    on th.forum = f.slug
			left join post p 
			    on p.forum = f.slug and p.threadid = th.id
		where f.slug = $1 group by f.slug
	`, forumParam).
		Scan(
			&forum.Slug,
			&forum.Title,
			&forum.User,
			&forum.Threads,
			&forum.Posts,
		)
	if err != nil {
		fmt.Println("CreateForum (2):", err)
		return eCtx.JSON(404, forms.Error{
			Message: "none such forum"})
	}

	return eCtx.JSON(200, forum)

}

func CreateForumThread(eCtx echo.Context) error {
	slug := eCtx.Param(forumSlug)
	thread := forms.Thread{}
	if err := eCtx.Bind(&thread); err != nil {
		fmt.Println("CreateForumThread (1):", err)
		return err
	}
	thread.Forum = sql.NullString{String: slug}

	ctx := eCtx.Request().Context()
	// try insert
	if err := DBS.QueryRowContext(ctx, `
		insert into thread (title,author,forum,message)
		values ($1,$2,$3,$4) returning id,created`,
		thread.Title,
		thread.Author,
		thread.Forum,
		thread.Message,
	).
		Scan(
			&thread.Id,
			&thread.Created,
		); err == nil {
		return eCtx.JSON(201, forms.ThreadResult{
			Id:      thread.Id,
			Title:   thread.Title,
			Author:  thread.Author,
			Forum:   thread.Forum.String,
			Message: thread.Message,
			Votes:   thread.Votes,
			Created: thread.Created,
		})
	}
	fmt.Println("forum:", thread.Forum, "author:", thread.Author)
	// select if exists
	if err := DBS.QueryRowContext(eCtx.Request().Context(), `
		select 
		th.id,
		th.title,
		th.author,
		th.forum,
		th.message,
		th.slug,
		th.created,
		count(v)
			from thread th 
			    left join vote v 
			        on v.threadid = th.id
		where th.forum = $1 and th.author = $2 
		group by th.id`,
		thread.Forum.String,
		thread.Author,
	).
		Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Created,
			&thread.Votes,
		); err == nil {
		return eCtx.JSON(409, forms.ThreadResult{
			Id:      thread.Id,
			Title:   thread.Title,
			Author:  thread.Author,
			Forum:   thread.Forum.String,
			Message: thread.Message,
			Votes:   thread.Votes,
			Created: thread.Created,
		})
	}
	fmt.Println("CreateForumThread not found (2)")
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
	fmt.Println("result: ", build.String())
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
	threadsFilter := forms.ThreadFilter{}
	if err := eCtx.Bind(&threadsFilter); err != nil {
		fmt.Println("GetForumUsers (1):", err)
		return err
	}
	build := strings.Builder{}
	build.WriteString(`
		select id,title,author,forum, message,
			votes, slug, created from thread
			where forum = $1`)

	if threadsFilter.Since != "" {
		build.WriteString(fmt.
			Sprintf(
				` and created >= %s`,
				threadsFilter.Since,
			),
		)
	}
	build.WriteString(`
		  ORDER BY created 
		    limit nullif($2,0)`)
	if threadsFilter.Desc {
		build.WriteString(" Desc")
	}
	if rows, err := DBS.
		QueryContext(
			ctx,
			build.String(),
			slug,
			threadsFilter.Limit); err == nil {

		threadsResponse := make([]forms.ThreadResult, 0, 100)
		for rows.Next() {
			thread := forms.Thread{}
			if err = rows.
				Scan(
					&thread.Id,
					&thread.Title,
					&thread.Author,
					&thread.Forum,
					&thread.Message,
					&thread.Votes,
					&thread.Slug,
					&thread.Created,
				); err != nil {
				fmt.Println("GetForumUsers (2):", err)

				return eCtx.JSON(http.StatusInternalServerError, forms.Error{
					Message: "smth wrong"})
			}
			threadsResponse = append(threadsResponse, forms.ThreadResult{
				Id:      thread.Id,
				Title:   thread.Title,
				Author:  thread.Author,
				Forum:   thread.Forum.String,
				Message: thread.Message,
				Votes:   thread.Votes,
				Created: thread.Created,
			})
		}
		return eCtx.JSON(http.StatusOK, threadsResponse)
	}
	return eCtx.JSON(http.StatusNotFound, forms.Error{
		Message: "none such forum"})
}
