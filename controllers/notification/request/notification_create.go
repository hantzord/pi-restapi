package request

type NotificationCreateRequest struct {
	DoctorID int    `json:"doctor_id"`
	UserID   int    `json:"user_id"`
	Content  string `json:"content" validator:"required"`
}
