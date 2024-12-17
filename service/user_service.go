package service

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"project/domain"
	"project/helper"
	"project/repository"
	"time"
)

type UserService interface {
	All(user domain.User) ([]domain.User, error)
	Get(user domain.User) (*domain.User, error)
	Register(user *domain.User) error
	UpdatePassword(id uuid.UUID, newPassword string) error
}

type userService struct {
	repo repository.Repository
	log  *zap.Logger
}

func NewUserService(repo repository.Repository, log *zap.Logger) UserService {
	return &userService{repo, log}
}

func (s *userService) All(user domain.User) ([]domain.User, error) {
	return s.repo.User.All(user)
}

func (s *userService) Get(user domain.User) (*domain.User, error) {
	return s.repo.User.Get(user)
}

func (s *userService) Register(user *domain.User) error {
	return s.repo.User.Create(user)
}

func (s *userService) UpdatePassword(id uuid.UUID, newPassword string) error {
	passwordResetToken := domain.PasswordResetToken{ID: id}
	if err := s.repo.PasswordReset.Get(&passwordResetToken); err != nil {
		return err
	}

	if passwordResetToken.ValidatedAt == nil {
		return errors.New("password reset token is invalid")
	}

	if passwordResetToken.PasswordResetAt != nil {
		return errors.New("password reset token has expired")
	}

	passwordResetToken.User.Password = helper.HashPassword(newPassword)
	if err := s.repo.User.Update(&passwordResetToken.User); err != nil {
		return err
	}

	passwordResetToken.PasswordResetAt = helper.Ptr(time.Now())
	if err := s.repo.PasswordReset.Update(&passwordResetToken); err != nil {
		return err
	}
	return nil
}
