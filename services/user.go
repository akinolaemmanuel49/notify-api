package services

import (
	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/repositories"
	"github.com/akinolaemmanuel49/notify-api/utils"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) CreateUser(userInputWithPassword *models.UserInputWithPassword) error {
	userInput := models.UserInput{
		FirstName: userInputWithPassword.FirstName,
		LastName:  userInputWithPassword.LastName,
		Email:     userInputWithPassword.Email,
	}

	password := userInputWithPassword.Password

	err := s.userRepository.CreateUser(&userInput, password)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) GetUserByID(id int64) (*models.UserProfile, error) {
	user, err := s.userRepository.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAllUsers(page, pageSize int) ([]*models.UserProfile, error) {
	users, err := s.userRepository.GetAllUsers(page, pageSize)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateUserByID(ID, claimsID int64, fields map[string]interface{}) error {
	if claimsID != ID {
		// utils.RespondWithError(w, "You are not permitted to modify this resource", http.StatusForbidden)
		return utils.ErrForbidden
	}
	err := s.userRepository.UpdateUserByID(ID, fields)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUserByID(ID, claimsID int64) error {
	if claimsID != ID {
		// utils.RespondWithError(w, "You are not permitted to modify this resource", http.StatusForbidden)
		return utils.ErrForbidden
	}
	err := s.userRepository.DeleteUserByID(ID)
	if err != nil {
		return err
	}
	return nil
}
