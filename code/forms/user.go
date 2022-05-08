package forms

type User struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
	Nickname string `json:"nickname,omitempty"`
}
