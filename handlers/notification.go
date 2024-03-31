package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/services"
	"github.com/akinolaemmanuel49/notify-api/utils"
	"github.com/golang-jwt/jwt/v5"
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
	log.Println("Endpoint Hit: Notifications HealthCheck")

	// Set response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write success response
	fmt.Fprintf(w, `{"status": "ok"}`)
}

func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: CreateNotification")

	// Extract user ID from JWT token
	token, err := utils.GetToken(r)
	if err != nil {
		utils.RespondWithError(w, "Error: unauthorized access", http.StatusUnauthorized)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	publisherID := int64(claims["id"].(float64))

	var notificationInput models.NotificationInput

	// Check and resolve errors during JSON decoding process
	err = json.NewDecoder(r.Body).Decode(&notificationInput)
	if err != nil {
		utils.RespondWithError(w, "Error: failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check and resolve errors from the create notification service
	err = h.notificationService.CreateNotification(&notificationInput, publisherID)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidRangeForPriority) {
			utils.RespondWithError(w, err.Error(), http.StatusBadRequest)
			return
		}
		utils.RespondWithError(w, fmt.Sprintf("Error: failed to create notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := models.NotificationResponse{
		Code:    http.StatusCreated,
		Message: "Notification was successfully created",
	}

	// Write response header
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *NotificationHandler) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: GetNotificationByID")

	vars := mux.Vars(r)

	// Convert string to integer
	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	// Check and resolve errors arising from string conversion
	if err != nil {
		utils.RespondWithError(w, "Error: invalid notification ID", http.StatusBadRequest)
		return
	}

	notification, err := h.notificationService.GetNotificationByID(ID)

	// Check and resolve errors from get notification by id service
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			utils.RespondWithError(w, fmt.Sprintf("Error: notification with id: %d was not found", ID), http.StatusNotFound)
			return
		}
		utils.RespondWithError(w, fmt.Sprintf("Error: failed to retrieve notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := models.NotificationResponse{
		Code:    http.StatusOK,
		Data:    notification,
		Message: fmt.Sprintf("Notification with ID: %d was successfully retrieved", ID),
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(response)
}

func (h *NotificationHandler) GetOwnNotifications(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: GetOwnNotifications")

	// Extract user ID from JWT token
	token, err := utils.GetToken(r)
	if err != nil {
		utils.RespondWithError(w, "Error: unauthorized access", http.StatusUnauthorized)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	publisherID := int64(claims["id"].(float64))

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

	notifications, err := h.notificationService.GetOwnNotifications(publisherID, page, pageSize)

	// Check and resolve errors from get all notifications service
	if err != nil {
		utils.RespondWithError(w, fmt.Sprintf("Error: failed to retrieve notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := models.NotificationResponse{
		Code:    http.StatusOK,
		Data:    notifications,
		Message: "Notifications successfully retrieved.",
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(response)
}

func (h *NotificationHandler) GetAllNotifications(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: GetAllNotifications")

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
		utils.RespondWithError(w, fmt.Sprintf("Error: failed to retrieve notification: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := models.NotificationResponse{
		Code:    http.StatusOK,
		Data:    notifications,
		Message: "Notifications successfully retrieved.",
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(response)
}

func (h *NotificationHandler) UpdateNotificationByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: UpdateNotificationByID")

	vars := mux.Vars(r)

	// Convert string to integer
	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	// Check and resolve errors arising from string conversion
	if err != nil {
		utils.RespondWithError(w, "Error: invalid notification ID", http.StatusBadRequest)
		return
	}

	// Extract user ID from JWT token
	token, err := utils.GetToken(r)
	if err != nil {
		utils.RespondWithError(w, "Error: unauthorized access", http.StatusUnauthorized)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	publisherID := int64(claims["id"].(float64))

	var fields map[string]interface{}

	// Check and resolve errors during JSON decoding process
	err = json.NewDecoder(r.Body).Decode(&fields)
	if err != nil {
		utils.RespondWithError(w, "Error: failed to parse request body", http.StatusBadRequest)
		return
	}

	err = h.notificationService.UpdateNotificationByID(ID, publisherID, fields)

	// Check and resolve errors from get notification by id service
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			utils.RespondWithError(w, fmt.Sprintf("Error: notification with id: %d was not found", ID), http.StatusNotFound)
			return
		}
		if errors.Is(err, utils.ErrInvalidTypeForPriority) {
			utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusBadRequest)
			return
		}
		if errors.Is(err, utils.ErrInvalidRangeForPriority) {
			utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusBadRequest)
			return
		}
		utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := models.NotificationResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Notification with ID: %d was successfully updated", ID),
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(response)
}

func (h *NotificationHandler) DeleteNotificationByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: DeleteNotificationByID")

	vars := mux.Vars(r)

	// Convert string to integer
	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	// Check and resolve errors arising from string conversion
	if err != nil {
		utils.RespondWithError(w, "Error: invalid notification ID", http.StatusBadRequest)
		return
	}

	// Extract user ID from JWT token
	token, err := utils.GetToken(r)
	if err != nil {
		utils.RespondWithError(w, "Error: unauthorized access", http.StatusUnauthorized)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	publisherID := int64(claims["id"].(float64))

	err = h.notificationService.DeleteNotificationByID(ID, publisherID)

	// Check and resolve errors from get notification by id service
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			utils.RespondWithError(w, fmt.Sprintf("Error: notification with id: %d was not found", ID), http.StatusNotFound)
			return
		}
		utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := models.NotificationResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Notification with ID: %d was successfully deleted", ID),
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(response)
}
