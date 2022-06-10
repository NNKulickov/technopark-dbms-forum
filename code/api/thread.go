package api

import (
	"context"
	"fmt"
	"github.com/NNKulickov/technopark-dbms-forum/forms"
	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func CreateThreadPost(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	thread, err := getThread(eCtx)
	if err != nil {
		fmt.Println("CreateThreadPost (1)", err)
		return err
	}

	posts := make([]forms.Post, 0, 100)
	if err = eCtx.Bind(&posts); err != nil {
		fmt.Println("CreateThreadPost (2)", err)
	}

	if len(posts) == 0 {
		return eCtx.JSON(http.StatusCreated, posts)
	}

	builder := strings.Builder{}

	builder.WriteString(`insert into post 
		(parent,author,message,isedited,forum,threadid,created) values`)

	args := []any{}

	//parentExists := false
	for i, post := range posts {
		//if post.Parent == 0 {
		//	parentExists = true
		//}
		builder.WriteString(fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,now()),",
			6*i+1, 6*i+2, 6*i+3, 6*i+4, 6*i+5, 6*i+6))
		post.IsEdited = false
		post.Forum = thread.Forum
		post.Thread = thread.Id
		args = append(args,
			post.Parent,
			post.Author,
			post.Message,
			post.IsEdited,
			post.Forum,
			post.Thread,
		)
	}
	//if !parentExists {
	//	fmt.Println("CreateThreadPost (3) None parent")
	//	return eCtx.JSON(http.StatusConflict, forms.Error{
	//		Message: "None parent",
	//	})
	//}

	sqlQuery := builder.String()
	sqlQuery = strings.TrimSuffix(sqlQuery, ",")
	sqlQuery += ` returning id,parent,author,message,
		isedited,forum,threadid,created`
	prep, err := DBS.PrepareContext(ctx, sqlQuery)
	if err != nil {
		fmt.Println("CreateThreadPost (4):", err)
		return err
	}
	rows, err := prep.QueryContext(ctx, args...)
	if err != nil {
		fmt.Println("CreateThreadPost (5):", err)
		return err
	}
	defer rows.Close()
	postsResult := make([]forms.Post, 0, 100)
	for rows.Next() {
		post := forms.Post{}
		err = rows.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if err != nil {
			fmt.Println("CreateThreadPost (6):", err)
			return err
		}
		postsResult = append(postsResult, post)
	}
	return eCtx.JSON(http.StatusCreated, postsResult)
}

func GetThreadDetails(eCtx echo.Context) error {

	//ctx := eCtx.Request().Context()
	thread, err := getThread(eCtx)
	if err != nil {
		fmt.Println("GetThreadDetails (1)", err)
		return err
	}

	return eCtx.JSON(http.StatusOK, forms.ThreadForm{
		Id:      thread.Id,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   thread.Forum,
		Message: thread.Message,
		Slug:    thread.Slug.String,
		Created: strfmt.DateTime(thread.Created.UTC()).String(),
	})
}
func UpdateThreadDetails(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	thread, err := getThread(eCtx)
	if err != nil {
		fmt.Println("UpdateThreadDetails (1)", err)
		return err
	}
	threadUpdate := forms.ThreadUpdate{}
	if err = eCtx.Bind(&threadUpdate); err != nil {
		fmt.Println("UpdateThreadDetails (2)", err)
	}
	if threadUpdate.Message != "" {
		thread.Message = threadUpdate.Message
	}
	if threadUpdate.Title != "" {
		thread.Title = threadUpdate.Title
	}
	if _, err = DBS.ExecContext(ctx, `
		update thread set title = $1, message = $2 where id = $3`, thread.Title, thread.Message, thread.Id); err != nil {
		fmt.Println("GetThreadDetails (3) none thread", err)
		return eCtx.JSON(http.StatusNotFound, forms.Error{
			Message: "None thread",
		})
	}
	return eCtx.JSON(http.StatusOK, forms.ThreadForm{
		Id:      thread.Id,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   thread.Forum,
		Message: thread.Message,
		Slug:    thread.Slug.String,
		Created: strfmt.DateTime(thread.Created.UTC()).String(),
	})
}

func SetThreadVote(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	thread, err := getThread(eCtx)
	if err != nil {
		fmt.Println("SetThreadVote (1)", err)
		return err
	}
	vote := forms.Vote{}
	if err = eCtx.Bind(&vote); err != nil {
		fmt.Println("SetThreadVote (2)", err)
	}
	if _, err = DBS.ExecContext(ctx, `
		insert into vote (threadid, nickname, voice)
			values ($1,$2,$3)
		on conflict on constraint unique_voice do update 
		set voice = excluded.voice`, thread.Id, vote.Nickname, vote.Voice); err != nil {
		fmt.Println("SetThreadVote (3)", err)
		return eCtx.JSON(http.StatusNotFound, forms.Error{
			Message: "None thread",
		})
	}
	votes, err := getVotes(ctx, thread.Id)
	if err != nil {
		fmt.Println("SetThreadVote (3)", err)
		return eCtx.JSON(http.StatusNotFound, forms.Error{
			Message: "None thread",
		})
	}
	return eCtx.JSON(http.StatusOK, forms.ThreadForm{
		Id:      thread.Id,
		Title:   thread.Title,
		Author:  thread.Author,
		Forum:   thread.Forum,
		Message: thread.Message,
		Slug:    thread.Slug.String,
		Created: strfmt.DateTime(thread.Created.UTC()).String(),
		Votes:   votes,
	})
}

func getThread(eCtx echo.Context) (forms.ThreadModel, error) {
	threadIdOrSlug, err := getThreadParam(eCtx)
	if err != nil {
		return forms.ThreadModel{}, err
	}
	ctx := eCtx.Request().Context()
	slug := ""
	id, err := strconv.Atoi(threadIdOrSlug)
	if err != nil {
		id = 0
		slug = threadIdOrSlug
	}
	thread := forms.ThreadModel{}
	if err = DBS.QueryRowContext(ctx, `
		select id, title, author, forum, message, slug, created from thread where id = $1 or lower(slug) = lower($2) `, id, slug).
		Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Created,
		); err != nil {
		fmt.Println("SetThreadVote (1) none thread", err)
		return forms.ThreadModel{}, eCtx.JSON(http.StatusNotFound, forms.Error{
			Message: "None thread",
		})
	}
	return thread, nil
}

func getThreadParam(eCtx echo.Context) (string, error) {
	slug := eCtx.Param(threadSlug)
	if slug == "" {
		eCtx.Logger().Error("cannot get username")
		return slug, eCtx.JSON(
			http.StatusBadRequest,
			forms.Error{Message: "None user"})
	}
	return slug, nil
}

func getVotes(ctx context.Context, threadid int) (int, error) {
	votes := 0
	err := DBS.QueryRowContext(ctx,
		`select sum(voice) from vote where threadid = $1`,
		threadid,
	).
		Scan(&votes)

	return votes, err
}

func GetThreadPosts(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	thread, err := getThread(eCtx)
	if err != nil {
		fmt.Println("GetThreadPosts (1)", err)
		return err
	}

	limitString := eCtx.QueryParam("limit")
	sinceString := eCtx.QueryParam("since")
	descString := eCtx.QueryParam("desc")
	limit := 0
	since := 0
	desc := false
	if limitString != "" {
		if limit, err = strconv.Atoi(limitString); err != nil {
			fmt.Println("GetThreadPosts (2)", err)
			return err
		}
	}
	if sinceString != "" {
		if since, err = strconv.Atoi(sinceString); err != nil {
			fmt.Println("GetThreadPosts (3)", err)
			return err
		}
	}
	if descString != "" {
		if desc, err = strconv.ParseBool(descString); err != nil {
			fmt.Println("GetThreadPosts (4)", err)
			return err
		}
	}
	postsMeta := forms.ThreadPosts{
		Limit: limit,
		Desc:  desc,
		Sort:  eCtx.QueryParam("sort"),
		Since: since,
	}
	if postsMeta.Sort == "" {
		postsMeta.Sort = "flat"
	}
	posts := []forms.Post{}
	switch postsMeta.Sort {
	case "flat":
		posts, err = getPostsFlat(ctx, thread.Id, postsMeta.Limit, postsMeta.Since, postsMeta.Desc)
	case "tree":
		posts, err = getPostsTree(ctx, thread.Id, postsMeta.Limit, postsMeta.Since, postsMeta.Desc)
	case "parent_tree":
		posts, err = getPostsParentTree(ctx, thread.Id, postsMeta.Limit, postsMeta.Since, postsMeta.Desc)
	}
	if err != nil {
		fmt.Println("GetThreadPosts (5)", err)
		return eCtx.JSON(http.StatusInternalServerError, forms.Error{Message: "parse err"})
	}

	return eCtx.JSON(http.StatusOK, posts)
}
