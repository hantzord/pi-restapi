package request

type SendFeedbackRequest struct {
	DoctorId uint   `json:"doctor_id" form:"doctor_id"`
	Rate     int    `json:"rate" form:"rate"`
	Message  string `json:"message" form:"message"`
}