package service

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"project/domain"
	"project/helper"
	"project/repository"
	"time"
)

type PasswordResetService interface {
	Create(token *domain.PasswordResetToken) error
	Validate(id uuid.UUID, token string) error
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

func (s *passwordResetService) Validate(id uuid.UUID, token string) error {
	passwordResetToken := domain.PasswordResetToken{ID: id, Otp: token, ValidatedAt: helper.Ptr(time.Now())}
	return s.repo.Update(&passwordResetToken)
}
