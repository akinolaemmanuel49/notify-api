package repositories

import (
	"database/sql"
	"errors"
	"log"

	"github.com/akinolaemmanuel49/notify-api/models"
)

var (
	ErrNotExists = errors.New("record does not exist")
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

// MigrateNotificationTable creates a notifications table if it does not exist.
func (r *NotificationRepository) MigrateNotificationTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS notifications(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		message TEXT NOT NULL,
		priority INTEGER NOT NULL DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := r.db.Exec(query)
	if err != nil {
		log.Panicln("Error migrating notification table:", err)
	}
	return err
}

// CreateNotification creates a new instance of NotificationRepository.
func (r *NotificationRepository) CreateNotification(notification *models.Notification) error {
	query := `
	INSERT INTO notifications(
		title, 
		message,
		priority,
		created_at) 
	VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, notification.Title, notification.Message, notification.Priority, notification.CreatedAt)
	if err != nil {
		log.Panicln("Error inserting notification:", err)
	}
	return err
}

// GetNotificationByID retrieves a notification by its ID from the database.
func (r *NotificationRepository) GetNotificationByID(id int64) (*models.Notification, error) {
	query := `
	SELECT id, title, message, priority, created_at FROM notifications WHERE id = ?`

	result := r.db.QueryRow(query, id)

	var notification models.Notification
	err := result.Scan(&notification.ID, &notification.Title, &notification.Message, &notification.Priority, &notification.CreatedAt)
	if err != nil {
		log.Panicln("Error retrieving notification:", err)
		return nil, err
	}
	return &notification, nil
}

// GetAllNotifications retrieves all notifications from the database with pagination.
func (r *NotificationRepository) GetAllNotifications(page, pageSize int) ([]*models.Notification, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	query := `
	SELECT id, title, message, priority, created_at FROM notifications LIMIT ? OFFSET ?`
	results, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		log.Panicln("Error retrieving notifications:", err)
		return nil, err
	}
	defer results.Close()

	notifications := []*models.Notification{}
	for results.Next() {
		var notification models.Notification
		err := results.Scan(&notification.ID, &notification.Title, &notification.Message, &notification.Priority, &notification.CreatedAt)
		if err != nil {
			log.Panicln("Error scanning notification row:", err)
			return nil, err
		}
		notifications = append(notifications, &notification)
	}
	if err := results.Err(); err != nil {
		log.Panicln("Error iterationg over notification rows:", err)
		return nil, err
	}
	return notifications, nil
}
