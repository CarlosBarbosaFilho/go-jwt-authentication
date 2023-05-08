package response

type UserResponse struct {
	Code    int64 `json:"code"`
	Message string
	Data    interface{}
}

type DataUser struct {
	ID       uint
	Name     string
	Email    string `gorm:"unique"`
	Password string
}
