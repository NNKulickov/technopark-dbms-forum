package forms

import (
	"database/sql"
	"time"
)

type ThreadFilter struct {
	Limit int    `json:"limit,omitempty"`
	Since string `json:"since,omitempty"`
	Desc  bool   `json:"desc,omitempty"`
}

type ThreadModel struct {
	Id      int
	Title   string
	Author  string
	Forum   string
	Message string
	Votes   int
	Slug    sql.NullString

	Created time.Time
}

type ThreadUpdate struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

type ThreadForm struct {
	Id      int    `json:"id,omitempty"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum,omitempty"`
	Message string `json:"message"`
	Votes   int    `json:"votes,omitempty"`
	Slug    string `json:"slug"`
	Created string `json:"created,omitempty"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}
