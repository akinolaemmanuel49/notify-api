package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) UserHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: Users HealthCheck")

	// Set response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write success response
	fmt.Fprint(w, `{"status": "ok"}`)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: CreateUser")

	var user models.User

	// Check and resolve errors during JSON decoding process
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check and resolve errors from the create user service
	err = h.userService.CreateUser(&user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Write response header
	w.WriteHeader(http.StatusCreated)
}
