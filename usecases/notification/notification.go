package notification

import (
	"capstone/entities"
	notificationEntities "capstone/entities/notification"
)

type NotificationUseCase struct {
	notificationRepository notificationEntities.NotificationRepository
}

func NewNotificationUseCase(notificationRepository notificationEntities.NotificationRepository) notificationEntities.NotificationUseCase {
	return &NotificationUseCase{notificationRepository}
}

func (usecase *NotificationUseCase) GetNotificationByUserID(userID int, metadata *entities.Metadata) (*[]notificationEntities.UserNotification, error) {
	notifications, err := usecase.notificationRepository.GetNotificationByUserID(userID, metadata)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (usecase *NotificationUseCase) CreateUserNotification(userID int, content string) error {
	notification := notificationEntities.ToUserNotification(uint(userID), content)
	if err := usecase.notificationRepository.CreateUserNotification(&notification); err != nil {
		return err
	}
	return nil
}

func (usecase *NotificationUseCase) DeleteUserNotification(notificationID int) error {
	if err := usecase.notificationRepository.DeleteUserNotification(notificationID); err != nil {
		return err
	}
	return nil
}

func (usecase *NotificationUseCase) UpdateStatusUserNotification(notificationID int) error {
	if err := usecase.notificationRepository.UpdateStatusUserNotification(notificationID); err != nil {
		return err
	}
	return nil
}

func (usecase *NotificationUseCase) GetNotificationByDoctorID(doctorID int, metadata *entities.Metadata) (*[]notificationEntities.DoctorNotification, error) {
	notifications, err := usecase.notificationRepository.GetNotificationByDoctorID(doctorID, metadata)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (usecase *NotificationUseCase) CreateDoctorNotification(doctorID uint, content string) error {
	notification := notificationEntities.ToDoctorNotification(doctorID, content)
	if err := usecase.notificationRepository.CreateDoctorNotification(&notification); err != nil {
		return err
	}
	return nil
}

func (usecase *NotificationUseCase) DeleteDoctorNotification(notificationID int) error {
	if err := usecase.notificationRepository.DeleteDoctorNotification(notificationID); err != nil {
		return err
	}
	return nil
}

func (usecase *NotificationUseCase) UpdateStatusDoctorNotification(notificationID int) error {
	if err := usecase.notificationRepository.UpdateStatusDoctorNotification(notificationID); err != nil {
		return err
	}
	return nil
}
