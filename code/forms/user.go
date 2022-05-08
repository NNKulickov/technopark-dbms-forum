package forms

type Profile struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

type User struct {
	Profile
	Nickname string `json:"nickname"`
}
