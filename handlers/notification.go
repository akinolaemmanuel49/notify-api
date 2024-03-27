package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/services"
	"github.com/akinolaemmanuel49/notify-api/utils"
	"github.com/gorilla/mux"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) NotificationHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Notifications HealthCheck")

	// Set response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write success response
	fmt.Fprintf(w, `{"status": "ok"}`)
}

func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: CreateNotification")

	var notification models.Notification

	// Check and resolve errors during JSON decoding process
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check and resolve errors from the create notification service
	err = h.notificationService.CreateNotification(&notification)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidRangeForPriority) {
			utils.RespondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to create notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Write response header
	w.WriteHeader(http.StatusCreated)
}

func (h *NotificationHandler) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: GetNotificationByID")

	vars := mux.Vars(r)

	// Convert string to integer
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	// Check and resolve errors arising from string conversion
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	notification, err := h.notificationService.GetNotificationByID(id)

	// Check and resolve errors from get notification by id service
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			http.Error(w, fmt.Sprintf("Notification with id: %d was not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to retrieve notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(notification)
}

func (h *NotificationHandler) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: GetAllNotifications")

	// Check the page query in the url, convert it to an integer, resolve errors
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	// Check the pageSize query in the url, convert it to an integer, resolve errors
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // default page size
	}

	notifications, err := h.notificationService.GetAllNotifications(page, pageSize)

	// Check and resolve errors from get all notifications service
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(notifications)
}

func (h *NotificationHandler) UpdateNotificationByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: UpdateNotificationByID")

	vars := mux.Vars(r)

	// Convert string to integer
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	// Check and resolve errors arising from string conversion
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	var fields map[string]interface{}

	// Check and resolve errors during JSON decoding process
	err = json.NewDecoder(r.Body).Decode(&fields)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	err = h.notificationService.UpdateNotificationByID(id, fields)

	// Check and resolve errors from get notification by id service
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			http.Error(w, fmt.Sprintf("Notification with id: %d was not found", id), http.StatusNotFound)
			return
		}
		if errors.Is(err, utils.ErrInvalidTypeForPriority) {
			utils.RespondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, utils.ErrInvalidRangeForPriority) {
			utils.RespondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to retrieve notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *NotificationHandler) DeleteNotificationByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: DeleteNotificationByID")

	vars := mux.Vars(r)

	// Convert string to integer
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	// Check and resolve errors arising from string conversion
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	err = h.notificationService.DeleteNotificationByID(id)

	// Check and resolve errors from get notification by id service
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			http.Error(w, fmt.Sprintf("Notification with id: %d was not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to retrieve notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
