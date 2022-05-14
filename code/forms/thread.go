package forms

import (
	"database/sql"
	"time"
)

type Thread struct {
	Id      int            `json:"id,omitempty"`
	Title   string         `json:"title"`
	Author  string         `json:"author"`
	Forum   sql.NullString `json:"forum,omitempty"`
	Message string         `json:"message"`
	Votes   int            `json:"votes,omitempty"`
	Slug    sql.NullString `json:"slug,omitempty"`
	Created time.Time      `json:"created,omitempty"`
}

type ThreadResult struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   int       `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}

type ThreadUpdate struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}
