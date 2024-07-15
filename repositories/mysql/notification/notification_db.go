package notification

import (
	"capstone/entities/notification"
	"gorm.io/gorm"
)

type UserNotification struct {
	gorm.Model
	UserID  uint   `gorm:"column:user_id"`
	Content string `gorm:"column:content"`
	IsRead  bool   `gorm:"column:is_read"`
}

type DoctorNotification struct {
	gorm.Model
	DoctorID uint   `gorm:"column:doctor_id"`
	Content  string `gorm:"column:content;type:text"`
	IsRead   bool   `gorm:"column:is_read;default:false;not null"`
}

func (n *UserNotification) ToUserEntities() *notification.UserNotification {
	return &notification.UserNotification{
		ID:        n.ID,
		UserID:    n.UserID,
		Content:   n.Content,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt.String(),
	}
}

func (n *DoctorNotification) ToDoctorEntities() *notification.DoctorNotification {
	return &notification.DoctorNotification{
		ID:        n.ID,
		DoctorID:  n.DoctorID,
		Content:   n.Content,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt.String(),
	}
}

func ToNotificationUserModel(notification *notification.UserNotification) *UserNotification {
	return &UserNotification{
		UserID:  notification.UserID,
		Content: notification.Content,
	}
}

func ToNotificationDoctorModel(notification *notification.DoctorNotification) *DoctorNotification {
	return &DoctorNotification{
		DoctorID: notification.DoctorID,
		Content:  notification.Content,
	}
}
