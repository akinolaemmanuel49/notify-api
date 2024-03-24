package services

import (
	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/repositories"
)

type NotificationService struct {
	notificationRepository *repositories.NotificationRepository
}

func NewNotificationService(notificationRepository *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepository: notificationRepository,
	}
}

func (s *NotificationService) CreateNotification(notification *models.Notification) error {
	err := s.notificationRepository.CreateNotification(notification)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) GetNotificationByID(id int64) (*models.Notification, error) {
	notification, err := s.notificationRepository.GetNotificationByID(id)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (s *NotificationService) GetAllNotifications(page, pageSize int) ([]*models.Notification, error) {
	notifications, err := s.notificationRepository.GetAllNotifications(page, pageSize)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (s *NotificationService) UpdateNotificationByID(id int64, fields map[string]interface{}) error {
	err := s.notificationRepository.UpdateNotificationByID(id, fields)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) DeleteNotificationByID(id int64) error {
	err := s.notificationRepository.DeleteNotificationByID(id)
	if err != nil {
		return err
	}
	return nil
}
