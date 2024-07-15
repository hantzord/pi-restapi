package request

type OTPVerifyRequest struct {
	Email string `json:"email" form:"email"`
	Code  string `json:"code" form:"code"`
}