package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/repositories"
	"github.com/akinolaemmanuel49/notify-api/services"
	"github.com/akinolaemmanuel49/notify-api/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: GenerateToken")

	var credentials models.AuthCredentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	ok, ID, err := h.authService.AuthenticateUser(&credentials)
	if err != nil {
		if errors.Is(err, repositories.ErrInvalidCredentials) {
			utils.RespondWithError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to authenticate user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if ok {
		token, err := utils.GenerateJWT(ID)
		if err != nil {
			errorMessage := fmt.Sprintf("Failed to generate token: %s", err.Error())
			utils.RespondWithError(w, errorMessage, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}
