package forms

type PostForum struct {
	Title string `json:"title"`
	User  string `json:"user"`
	Slug  string `json:"slug"`
}

type ForumResult struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}
