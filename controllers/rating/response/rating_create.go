package response

type RatingCreateResponse struct {
	Id       uint   `json:"id"`
	UserId   uint   `json:"user_id"`
	DoctorId uint   `json:"doctor_id"`
	Rate     int    `json:"rate"`
	Message  string `json:"message"`
}