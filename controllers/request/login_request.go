package request

type UserRequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
