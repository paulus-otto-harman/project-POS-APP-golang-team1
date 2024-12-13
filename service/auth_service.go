package service

import (
	"go.uber.org/zap"
	"project/domain"
	"project/repository"
)

type AuthService interface {
	Login(user domain.User) (string, bool, error)
}

type authService struct {
	repo repository.AuthRepository
	log  *zap.Logger
}

func NewAuthService(repo repository.AuthRepository, log *zap.Logger) AuthService {
	return &authService{repo, log}
}

func (s *authService) Login(user domain.User) (string, bool, error) {
	return s.repo.Authenticate(user)
}
