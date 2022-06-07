package api

import (
	"context"
	"fmt"
	"github.com/NNKulickov/technopark-dbms-forum/forms"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func CreateThreadPost(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	thread, err := getThread(eCtx)
	if err != nil {
		fmt.Println("GetThreadDetails (1)", err)
		return err
	}

	posts := make([]forms.Post, 0, 100)
	if err = eCtx.Bind(&posts); err != nil {
		fmt.Println("CreateThreadPost (1)", err)
	}

	builder := strings.Builder{}

	builder.WriteString(`insert into post 
		(parent,author,message,isedited,forum,threadid,created) values`)

	args := []any{}

	parentExists := false
	for i, post := range posts {
		if post.Parent == 0 {
			parentExists = true
		}
		builder.WriteString(fmt.Sprintf("($%d,$%s,$%s,$%b,$%s,$%d,now()),",
			6*i+1, 6*i+2, 6*i+3, 6*i+4, 6*i+5, 6*i+6))
		post.IsEdited = false
		post.Forum = thread.Forum.String
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
	if !parentExists {
		fmt.Println("CreateThreadPost (2) None parent")
		return eCtx.JSON(http.StatusConflict, forms.Error{
			Message: "None parent",
		})
	}

	sqlQuery := builder.String()
	sqlQuery = strings.TrimSuffix(sqlQuery, ",")
	prep, err := DBS.PrepareContext(ctx, sqlQuery)
	if err != nil {
		fmt.Println("CreateThreadPost (3):", err)
		return err
	}
	if _, err = prep.ExecContext(ctx, args...); err != nil {
		fmt.Println("CreateThreadPost (4):", err)
		return err
	}
	return eCtx.JSON(http.StatusCreated, posts)
}

func GetThreadDetails(eCtx echo.Context) error {

	ctx := eCtx.Request().Context()
	thread, err := getThread(eCtx)
	if err != nil {
		fmt.Println("GetThreadDetails (1)", err)
		return err
	}

	votes, err := getVotes(ctx, thread.Id)
	if err != nil {
		fmt.Println("GetThreadDetails (2)", err)
		return eCtx.JSON(http.StatusNotFound, forms.Error{
			Message: "None thread",
		})
	}

	return eCtx.JSON(http.StatusOK, forms.ThreadResult{
		Id:      thread.Id,
		Title:   thread.Title,
		Forum:   thread.Forum.String,
		Message: thread.Message,
		Slug:    thread.Slug.String,
		Created: thread.Created,
		Votes:   votes,
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
	votes, err := getVotes(ctx, thread.Id)
	if err != nil {
		fmt.Println("GetThreadDetails (4)", err)
		return eCtx.JSON(http.StatusNotFound, forms.Error{
			Message: "None thread",
		})
	}
	return eCtx.JSON(http.StatusOK, forms.ThreadResult{
		Id:      thread.Id,
		Title:   thread.Title,
		Forum:   thread.Forum.String,
		Message: thread.Message,
		Slug:    thread.Slug.String,
		Created: thread.Created,
		Votes:   votes,
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
		insert into vote (threadid, nickname, voice) values ($1,$2,$3)`, thread.Id, vote.Nickname, vote.Voice); err != nil {
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

	return eCtx.JSON(http.StatusOK, forms.ThreadResult{
		Id:      thread.Id,
		Title:   thread.Title,
		Forum:   thread.Forum.String,
		Message: thread.Message,
		Slug:    thread.Slug.String,
		Created: thread.Created,
		Votes:   votes,
	})
}

func getThread(eCtx echo.Context) (forms.Thread, error) {
	threadIdOrSlug, err := getThreadParam(eCtx)
	if err != nil {
		return forms.Thread{}, err
	}
	ctx := eCtx.Request().Context()
	slug := ""
	id, err := strconv.Atoi(threadIdOrSlug)
	if err != nil {
		id = 0
		slug = threadIdOrSlug
	}
	thread := forms.Thread{}
	if err = DBS.QueryRowContext(ctx, `
		select id, title, author, forum, message, slug, created from thread where id = $1 or slug = $2 `, id, slug).
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
		return forms.Thread{}, eCtx.JSON(http.StatusNotFound, forms.Error{
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
	postsMeta := forms.ThreadPosts{}
	if err = eCtx.Bind(&postsMeta); err != nil {
		fmt.Println("GetThreadPosts (2)", err)
		return err
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
		fmt.Println("GetThreadPosts (3)", err)
		return eCtx.JSON(http.StatusInternalServerError, forms.Error{Message: "parse err"})
	}

	return eCtx.JSON(http.StatusOK, posts)
}
