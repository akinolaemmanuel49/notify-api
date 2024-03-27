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

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: GetUserByID")

	vars := mux.Vars(r)

	// Convert string to integer
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if errors.Is(err, utils.ErrNotFound) {
		http.Error(w, fmt.Sprintf("User with id: %d was not found", id), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: GetAllUsers")

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10 // default page size
	}

	users, err := h.userService.GetAllUsers(page, pageSize)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve users: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: UpdateUserByID")

	vars := mux.Vars(r)

	// Convert string to integer
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	// Check and resolve errors arising from string conversion
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var fields map[string]interface{}

	// Check and resolve errors during JSON decoding process
	err = json.NewDecoder(r.Body).Decode(&fields)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	err = h.userService.UpdateUserByID(id, fields)

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

func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: DeleteUserByID")

	vars := mux.Vars(r)

	// Convert string to integer
	id, err := strconv.ParseInt(vars["id"], 10, 64)

	// Check and resolve errors arising from string conversion
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUserByID(id)

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
