package response

type NotificationUserResponse struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type NotificationDoctorResponse struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}
