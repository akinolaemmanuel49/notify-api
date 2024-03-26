package services

import (
	"github.com/akinolaemmanuel49/notify-api/models"
	"github.com/akinolaemmanuel49/notify-api/repositories"
)

type AuthService struct {
	authRepository *repositories.AuthRepository
}

func NewAuthService(authRepository *repositories.AuthRepository) *AuthService {
	return &AuthService{
		authRepository: authRepository,
	}
}

func (s *AuthService) AuthenticateUser(credentials *models.AuthCredentials) (bool, int64, error) {
	ok, ID, err := s.authRepository.AuthenticateUser(credentials)
	if err != nil {
		return false, 0, err
	}
	return ok, ID, nil
}
