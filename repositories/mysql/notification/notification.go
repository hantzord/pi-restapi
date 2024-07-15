package notification

import (
	"capstone/constants"
	"capstone/entities"
	notificationEntities "capstone/entities/notification"
	"fmt"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) notificationEntities.NotificationRepository {
	return &NotificationRepository{db}
}

func (repository *NotificationRepository) GetNotificationByUserID(userID int, metadata *entities.Metadata) (*[]notificationEntities.UserNotification, error) {
	var notifications []UserNotification
	if err := repository.db.Limit(metadata.Limit).Offset(metadata.Offset()).Find(&notifications, "user_id = ?", userID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	var userNotifications []notificationEntities.UserNotification
	for _, notification := range notifications {
		userNotifications = append(userNotifications, *notification.ToUserEntities())
	}
	return &userNotifications, nil
}

func (repository *NotificationRepository) CreateUserNotification(notification *notificationEntities.UserNotification) error {
	notificationDB := ToNotificationUserModel(notification)
	if err := repository.db.Create(&notificationDB).Error; err != nil {
		return constants.ErrInsertDatabase
	}
	return nil
}

func (repository *NotificationRepository) DeleteUserNotification(notificationID int) error {
	if err := repository.db.Delete(&UserNotification{}, "id = ?", notificationID).Error; err != nil {
		return constants.ErrDeleteDatabase
	}
	return nil
}

func (repository *NotificationRepository) UpdateStatusUserNotification(notificationID int) error {
	var notification UserNotification
	if err := repository.db.First(&notification, "id = ?", notificationID).Error; err != nil {
		return constants.ErrDataNotFound
	}
	if notification.IsRead == true {
		return constants.ErrNotificationAlreadyRead
	}
	if err := repository.db.Model(&UserNotification{}).Where("id LIKE ?", notificationID).Update("is_read", true).Error; err != nil {
		return constants.ErrUpdateDatabase
	}
	return nil
}

func (repository *NotificationRepository) GetNotificationByDoctorID(doctorID int, metadata *entities.Metadata) (*[]notificationEntities.DoctorNotification, error) {
	var notifications []DoctorNotification
	if err := repository.db.Limit(metadata.Limit).Offset(metadata.Offset()).Find(&notifications, "doctor_id = ?", doctorID).Error; err != nil {
		return nil, constants.ErrDataNotFound
	}
	var doctorNotifications []notificationEntities.DoctorNotification
	for _, notification := range notifications {
		doctorNotifications = append(doctorNotifications, *notification.ToDoctorEntities())
	}
	return &doctorNotifications, nil
}

func (repository *NotificationRepository) CreateDoctorNotification(notification *notificationEntities.DoctorNotification) error {
	notificationDB := ToNotificationDoctorModel(notification)
	if err := repository.db.Create(&notificationDB).Error; err != nil {
		return constants.ErrInsertDatabase
	}
	return nil
}

func (repository *NotificationRepository) DeleteDoctorNotification(notificationID int) error {
	if err := repository.db.Delete(&DoctorNotification{}, "id = ?", notificationID).Error; err != nil {
		return constants.ErrDeleteDatabase
	}
	return nil
}

func (repository *NotificationRepository) UpdateStatusDoctorNotification(notificationID int) error {
	var notification DoctorNotification
	fmt.Println(notificationID)
	if err := repository.db.First(&notification, "id = ?", notificationID).Error; err != nil {
		return constants.ErrDataNotFound
	}
	if notification.IsRead == true {
		return constants.ErrNotificationAlreadyRead
	}
	if err := repository.db.Model(&DoctorNotification{}).Where("id LIKE ?", notificationID).Update("is_read", true).Error; err != nil {
		return constants.ErrUpdateDatabase
	}
	return nil
}
