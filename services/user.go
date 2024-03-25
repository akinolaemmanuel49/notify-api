package services

import (
	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/repositories"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) CreateUser(user *models.User) error {
	err := s.userRepository.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}
