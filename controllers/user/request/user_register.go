package request

type UserRegisterRequest struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}