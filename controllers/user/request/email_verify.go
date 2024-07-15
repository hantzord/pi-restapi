package request

type VerifyOTPRequest struct {
	Code string `json:"code" form:"code"`
}
