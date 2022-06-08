package forms

type Post struct {
	Id       int    `json:"id,omitempty"`
	Parent   int    `json:"parent"`
	Author   string `json:"author"`
	Message  string `json:"message"`
	IsEdited bool   `json:"isEdited,omitempty"`
	Forum    string `json:"forum,omitempty"`
	Thread   int    `json:"thread"`
	Created  string `json:"created,omitempty"`
}

type PostFull struct {
	Post   Post       `json:"post"`
	Author User       `json:"author"`
	Thread ThreadForm `json:"thread"`
}

type ThreadPosts struct {
	Limit int    `json:"limit,omitempty"`
	Since int    `json:"since,omitempty"`
	Sort  string `json:"sort,omitempty"`
	Desc  bool   `json:"desc"`
}
