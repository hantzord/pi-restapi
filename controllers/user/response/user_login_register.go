package response

type UserLoginRegisterResponse struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
}