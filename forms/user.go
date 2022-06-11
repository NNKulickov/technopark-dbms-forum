package forms

type User struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
	Nickname string `json:"nickname,omitempty"`
}

type UserFilter struct {
	Limit int    `json:"limit,omitempty"`
	Since string `json:"since,omitempty"`
	Desc  bool   `json:"desc,omitempty"`
}
