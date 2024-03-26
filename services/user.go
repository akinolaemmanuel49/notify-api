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

func (s *UserService) GetUserByID(id int64) (*models.User, error) {
	user, err := s.userRepository.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAllUsers(page, pageSize int) ([]*models.User, error) {
	users, err := s.userRepository.GetAllUsers(page, pageSize)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateUserByID(id int64, fields map[string]interface{}) error {
	err := s.userRepository.UpdateUserByID(id, fields)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUserByID(id int64) error {
	err := s.userRepository.DeleteUserByID(id)
	if err != nil {
		return err
	}
	return nil
}
