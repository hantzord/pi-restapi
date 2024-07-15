package request

type OTPRequest struct {
	Email string `json:"email" form:"email"`
}