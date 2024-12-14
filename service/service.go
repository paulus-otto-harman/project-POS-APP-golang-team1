package service

import (
	"go.uber.org/zap"
	"project/repository"
)

type Service struct {
	Auth          AuthService
	PasswordReset PasswordResetService
	User          UserService
}

func NewService(repo repository.Repository, log *zap.Logger) Service {
	return Service{
		Auth:          NewAuthService(repo.User, log),
		PasswordReset: NewPasswordResetService(repo.PasswordReset, log),
		User:          NewUserService(repo.User, log),
	}
}
