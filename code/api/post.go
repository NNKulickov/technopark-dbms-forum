package api

import (
	"context"
	"fmt"
	"github.com/NNKulickov/technopark-dbms-forum/forms"
	"github.com/labstack/echo/v4"
	"strings"
)

func GetPostDetails(eCtx echo.Context) error    { return nil }
func UpdatePostDetails(eCtx echo.Context) error { return nil }

func getPostsFlat(ctx context.Context, threadid, limit, since int, desc bool) ([]forms.Post, error) {
	posts := []forms.Post{}
	builder := strings.Builder{}

	builder.WriteString(`
		select id,parent,author,message,isedited,forum,threadid,created
			from post where threadid = $1`)
	if since > 0 {
		builder.WriteString(fmt.Sprintf(" and id > %d", since))
	}

	builder.WriteString(" order by created ")

	if desc {
		builder.WriteString("desc")
	}
	builder.WriteString(" limit nullif($2,0)")
	rows, err := DBS.QueryContext(ctx, builder.String(),
		threadid, limit)
	if err != nil {
		fmt.Println("getPostsFlat(1): ", err)
		return nil, err
	}
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
			fmt.Println("getPostsFlat(2): ", err, post)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
func getPostsParentTree(ctx context.Context, threadid, limit, since int, desc bool) ([]forms.Post, error) {
	posts := []forms.Post{}
	builder := strings.Builder{}
	builder.WriteString(`with parents(id,parent,author,message,isedited,forum,threadid,created,path) as (`)
	builder.WriteString(`
		select id,parent,author,message,isedited,forum,threadid,created,path
			from post where threadid = $1 and parent = 0`)
	if since > 0 {
		builder.WriteString(fmt.Sprintf(" and id > %d", since))
	}

	builder.WriteString(" order by id")

	if desc {
		builder.WriteString("desc")
	}
	builder.WriteString(" limit nullif($2,0) )")
	builder.WriteString(`
				select id,parent,author,message,isedited,forum,threadid,created,path
			from post where threadid = $1 and id in parents.path
		 `)
	fmt.Println("getPostsParentTree:", builder.String())
	rows, err := DBS.QueryContext(ctx, builder.String(),
		threadid, limit)
	if err != nil {
		fmt.Println("getPostsParentTree(1): ", err)
		return nil, err
	}
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
			fmt.Println("getPostsParentTree(2): ", err, post)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPostsTree(ctx context.Context, threadid, limit, since int, desc bool) ([]forms.Post, error) {
	posts := []forms.Post{}
	builder := strings.Builder{}

	builder.WriteString(`
		select id,parent,author,message,isedited,forum,threadid,created
			from post where threadid = $1`)
	if since > 0 {
		builder.WriteString(fmt.Sprintf(" and id > %d", since))
	}

	builder.WriteString(" order by id,parent ")

	if desc {
		builder.WriteString("desc")
	}
	builder.WriteString(" limit nullif($2,0)")
	rows, err := DBS.QueryContext(ctx, builder.String(),
		threadid, limit)
	if err != nil {
		fmt.Println("getPostsTree(1): ", err)
		return nil, err
	}
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
			fmt.Println("getPostsTree(2): ", err, post)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
