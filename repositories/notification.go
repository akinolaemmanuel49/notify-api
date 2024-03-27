package repositories

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/utils"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

// CreateNotification creates a new instance of NotificationRepository.
func (r *NotificationRepository) CreateNotification(notification *models.Notification) error {
	query := `
	INSERT INTO notifications(
		title, 
		message,
		priority,
		publisher_id,
		created_at,
		updated_at) 
	VALUES (($1), ($2), ($3), ($4), ($5))`

	_, err := r.db.Exec(query,
		notification.Title,
		notification.Message,
		notification.Priority,
		notification.PublisherID,
		notification.CreatedAt,
		notification.UpdatedAt)
	if err != nil {
		log.Panicln("Error inserting notification:", err)
	}
	return err
}

// GetNotificationByID retrieves a notification by its ID from the database.
func (r *NotificationRepository) GetNotificationByID(id int64) (*models.Notification, error) {
	query := `
	SELECT 
		id, 
		title, 
		message, 
		priority, 
		publisher_id, 
		created_at, 
		updated_at 
	FROM notifications WHERE id = ($1)`

	result := r.db.QueryRow(query, id)

	var notification models.Notification
	err := result.Scan(&notification.ID, &notification.Title, &notification.Message, &notification.Priority, &notification.PublisherID, &notification.CreatedAt, &notification.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
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
	SELECT id, title, message, priority, publisher_id, created_at, updated_at FROM notifications LIMIT $1 OFFSET $2`
	results, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		log.Panicln("Error retrieving notifications:", err)
		return nil, err
	}
	defer results.Close()

	notifications := []*models.Notification{}
	for results.Next() {
		var notification models.Notification
		err := results.Scan(&notification.ID, &notification.Title, &notification.Message, &notification.Priority, &notification.CreatedAt, &notification.UpdatedAt)
		if err != nil {
			log.Panicln("Error scanning notification row:", err)
			return nil, err
		}
		notifications = append(notifications, &notification)
	}
	if err := results.Err(); err != nil {
		log.Panicln("Error iterating over notification rows:", err)
		return nil, err
	}
	return notifications, nil
}

func (r *NotificationRepository) UpdateNotificationByID(id int64, fields map[string]interface{}) error {
	_, err := r.GetNotificationByID(id)
	if errors.Is(err, utils.ErrNotFound) {
		return utils.ErrNotFound
	}

	query := "UPDATE notifications SET "

	var params []interface{}
	i := 1
	for key, value := range fields {
		if i > 1 {
			query += ", "
		}
		query += key + " = $" + strconv.Itoa(i)
		params = append(params, value)
		i++
	}

	query += " WHERE id = $" + strconv.Itoa(i)
	params = append(params, id)

	_, err = r.db.Exec(query, params...)
	if err != nil {
		log.Panicln("Error updating notification: ", err)
	}
	return err
}

func (r *NotificationRepository) DeleteNotificationByID(id int64) error {
	_, err := r.GetNotificationByID(id)
	if errors.Is(err, utils.ErrNotFound) {
		return utils.ErrNotFound
	}

	query := `
	DELETE FROM notifications WHERE id = ($1)`

	_, err = r.db.Exec(query, id)

	if err != nil {
		log.Panicln("Error deleting notification: ", err)
	}
	return err
}
