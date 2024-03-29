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
	"github.com/golang-jwt/jwt/v5"
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

	var userInputWithPassword models.UserInputWithPassword

	// Check and resolve errors during JSON decoding process
	err := json.NewDecoder(r.Body).Decode(&userInputWithPassword)
	if err != nil {
		utils.RespondWithError(w, "Error: failed to parse request body", http.StatusBadRequest)
		return
	}

	// Check and resolve errors from the create user service
	err = h.userService.CreateUser(&userInputWithPassword)
	if errors.Is(err, utils.ErrDuplicateKey) {
		utils.RespondWithError(w, "Error: email address already in use", http.StatusConflict)
		return
	}
	if err != nil {
		utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := models.UserResponse{
		Code:    http.StatusCreated,
		Message: "User was successfully created",
	}

	// Write response header
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: GetUserByID")

	vars := mux.Vars(r)

	// Convert string to integer
	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		utils.RespondWithError(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(ID)
	if errors.Is(err, utils.ErrNotFound) {
		utils.RespondWithError(w, fmt.Sprintf("Error: user with ID: %d was not found", ID), http.StatusNotFound)
		return
	}
	if err != nil {
		utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response := models.UserResponse{
		Code:    http.StatusOK,
		Data:    user,
		Message: fmt.Sprintf("User with ID: %d was successfully retrieved", ID),
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(response)
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
		utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := models.UserResponse{
		Code:    http.StatusOK,
		Data:    users,
		Message: "Users successfully retrieved",
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: UpdateUserByID")

	vars := mux.Vars(r)

	// Convert string to integer
	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	// Check and resolve errors arising from string conversion
	if err != nil {
		utils.RespondWithError(w, "Error: invalid user ID", http.StatusBadRequest)
		return
	}

	// Extract user ID from JWT token
	token, err := utils.GetToken(r)
	if err != nil {
		utils.RespondWithError(w, "Error: unauthorized access", http.StatusUnauthorized)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	claimsID := int64(claims["id"].(float64))

	var fields map[string]interface{}

	// Check and resolve errors during JSON decoding process
	err = json.NewDecoder(r.Body).Decode(&fields)
	if err != nil {
		utils.RespondWithError(w, "Error: failed to parse request body", http.StatusBadRequest)
		return
	}

	err = h.userService.UpdateUserByID(ID, claimsID, fields)

	// Check and resolve errors from get notification by ID service
	if err != nil {
		if errors.Is(err, utils.ErrForbidden) {
			utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusForbidden)
			return
		}
		if errors.Is(err, utils.ErrNotFound) {
			utils.RespondWithError(w, fmt.Sprintf("Error: notification with ID: %d was not found", ID), http.StatusNotFound)
			return
		}
		utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := models.UserResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("User with ID: %d was successfully updated", ID),
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: DeleteUserByID")

	vars := mux.Vars(r)

	// Convert string to integer
	ID, err := strconv.ParseInt(vars["id"], 10, 64)

	// Check and resolve errors arising from string conversion
	if err != nil {
		utils.RespondWithError(w, "Error: invalid user ID", http.StatusBadRequest)
		return
	}

	// Extract user ID from JWT token
	token, err := utils.GetToken(r)
	if err != nil {
		utils.RespondWithError(w, "Error: unauthorized access", http.StatusUnauthorized)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	claimsID := int64(claims["id"].(float64))

	err = h.userService.DeleteUserByID(ID, claimsID)

	// Check and resolve errors from get notification by id service
	if err != nil {
		if errors.Is(err, utils.ErrForbidden) {
			utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusForbidden)
			return
		}
		if errors.Is(err, utils.ErrNotFound) {
			utils.RespondWithError(w, fmt.Sprintf("Error: notification with ID: %d was not found", ID), http.StatusNotFound)
			return
		}
		utils.RespondWithError(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := models.UserResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("User with ID: %d was successfully deleted", ID),
	}

	// Encode and write JSON response
	json.NewEncoder(w).Encode(response)
}
