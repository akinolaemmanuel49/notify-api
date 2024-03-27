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

func (s *NotificationService) CreateNotification(notificationInput *models.NotificationInput, publisherID int64) error {
	if err := notificationInput.Priority.Validate(); err != nil {
		return err
	}
	err := s.notificationRepository.CreateNotification(notificationInput, publisherID)
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

func (s *NotificationService) UpdateNotificationByID(ID, publisherID int64, fields map[string]interface{}) error {
	if priorityField, ok := fields["priority"]; ok {
		// Validate priority
		priority, err := models.PriorityFromField(priorityField)
		if err != nil {
			return err
		}
		if err := priority.Validate(); err != nil {
			return err
		}
		fields["priority"] = priority
	}
	err := s.notificationRepository.UpdateNotificationByID(ID, publisherID, fields)
	if err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) DeleteNotificationByID(ID, publisherID int64) error {
	err := s.notificationRepository.DeleteNotificationByID(ID, publisherID)
	if err != nil {
		return err
	}
	return nil
}
