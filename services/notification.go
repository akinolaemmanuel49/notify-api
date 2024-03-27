package services

import (
	"fmt"
	"reflect"

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
	if err := notification.Priority.Validate(); err != nil {
		return err
	}
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
	fmt.Println(reflect.TypeOf(fields["priority"]))
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
