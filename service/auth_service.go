package service

import (
	"errors"
	"go.uber.org/zap"
	"project/domain"
	"project/helper"
	"project/repository"
)

type AuthService interface {
	Login(email, password string) (*domain.User, error)
}

type authService struct {
	repo repository.UserRepository
	log  *zap.Logger
}

func NewAuthService(repo repository.UserRepository, log *zap.Logger) AuthService {
	return &authService{repo, log}
}

func (s *authService) Login(email, password string) (*domain.User, error) {
	s.log.Info("Attempting to log in user", zap.String("email", email))

	// Cari user berdasarkan email
	user, err := s.repo.Get(domain.User{Email: email})
	if err != nil {
		s.log.Error("Login failed", zap.Error(err))
		return nil, err
	}

	// Verifikasi password
	if !helper.CheckPassword(password, user.Password) {
		s.log.Warn("Invalid password", zap.String("email", email))
		return nil, errors.New("invalid email or password")
	}

	s.log.Info("User logged in successfully", zap.String("email", email))
	return user, nil
}
