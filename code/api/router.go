package api

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var DBS *sql.DB

const (
	forumSlug    = "slug"
	postSlug     = "id"
	threadSlug   = "thread_slug"
	usernameSlug = "username"
)

func initForum(e *echo.Group) {
	const (
		create        = "/create"
		slugged       = "/:" + forumSlug
		details       = slugged + "/details"
		sluggedCreate = slugged + "/create"
		users         = slugged + "/users"
		threads       = slugged + "/threads"
	)
	e.POST(create, CreateForum)
	e.GET(details, GetForumDetails)
	e.POST(sluggedCreate, CreateForumThread)
	e.GET(users, GetForumUsers)
	e.GET(threads, GetForumThreads)
}

func initPost(e *echo.Group) {
	const (
		postDetails = "/:" + postSlug + "details"
	)
	e.GET(postDetails, GetPostDetails)
	e.POST(postDetails, UpdatePostDetails)
}

func initService(e *echo.Group) {
	const (
		clear  = "/clear"
		status = "/status"
	)
	e.GET(status, GetServiceStatus)
	e.POST(clear, ClearServiceData)
}

func initThread(e *echo.Group) {
	const (
		slugged = "/:" + threadSlug
		create  = slugged + "/create"
		details = slugged + "/details"
		posts   = slugged + "/posts"
		vote    = slugged + "/vote"
	)
	e.POST(create, CreateThreadPost)
	e.GET(details, GetThreadDetails)
	e.POST(details, UpdateThreadDetails)
	e.GET(posts, GetThreadPosts)
	e.POST(vote, SetThreadVote)

}

func initUser(e *echo.Group) {
	const (
		slugged = "/:" + usernameSlug
		create  = slugged + "/create"
		profile = slugged + "/profile"
	)
	e.POST(create, CreateUser)
	e.GET(profile, GetUserProfile)
	e.POST(profile, UpdateUserProfile)

}

func InitRoutes(e *echo.Group) {
	const (
		forum   = "/forum"
		post    = "/post"
		service = "/service"
		thread  = "/thread"
		user    = "/user"
	)
	initForum(e.Group(forum))
	initPost(e.Group(post))
	initService(e.Group(service))
	initThread(e.Group(thread))
	initUser(e.Group(user))
}
