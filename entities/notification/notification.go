package notification

import "capstone/controllers/notification/response"

type UserNotification struct {
	ID        uint
	UserID    uint
	Content   string
	IsRead    bool
	CreatedAt string
}

type DoctorNotification struct {
	ID        uint
	DoctorID  uint
	Content   string
	IsRead    bool
	CreatedAt string
}

func (n *UserNotification) ToUserResponse() *response.NotificationUserResponse {
	return &response.NotificationUserResponse{
		ID:        n.ID,
		Content:   n.Content,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}
}

func (n *DoctorNotification) ToDoctorResponse() *response.NotificationDoctorResponse {
	return &response.NotificationDoctorResponse{
		ID:        n.ID,
		Content:   n.Content,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
	}
}

func ToDoctorNotification(doctorID uint, content string) DoctorNotification {
	return DoctorNotification{
		DoctorID: doctorID,
		Content:  content,
	}
}

func ToUserNotification(userID uint, content string) UserNotification {
	return UserNotification{
		UserID:  userID,
		Content: content,
	}
}
