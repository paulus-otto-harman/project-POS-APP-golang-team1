package service

import (
	"go.uber.org/zap"
	"project/domain"
	"project/repository"
)

type PasswordResetService interface {
	Create(token *domain.PasswordResetToken) error
}

type passwordResetService struct {
	repo repository.PasswordResetRepository
	log  *zap.Logger
}

func NewPasswordResetService(repo repository.PasswordResetRepository, log *zap.Logger) PasswordResetService {
	return &passwordResetService{repo, log}
}

func (s *passwordResetService) Create(token *domain.PasswordResetToken) error {
	return s.repo.Create(token)
}
