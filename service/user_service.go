package service

import (
	"errors"
	"project/domain"
	"project/helper"
	"project/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService interface {
	All(sortField, sortDirection string, page, limit uint) ([]domain.User, int64, error)
	Get(user domain.User) (*domain.User, error)
	Register(user *domain.User) error
	UpdatePassword(id uuid.UUID, newPassword string) error
	Delete(id uint) error
	Update(user domain.User) error
}

type userService struct {
	repo repository.Repository
	log  *zap.Logger
}

func NewUserService(repo repository.Repository, log *zap.Logger) UserService {
	return &userService{repo, log}
}

func (s *userService) All(sortField, sortDirection string, page, limit uint) ([]domain.User, int64, error) {
	return s.repo.User.All(sortField, sortDirection, page, limit)
}

func (s *userService) Register(user *domain.User) error {
	existedUser := s.repo.User.GetByEmail(user.Email)
	if existedUser != nil {
		return errors.New("user email already exists")
	}

	if user.Role == "admin" {
		user.Password = helper.HashPassword(user.Password)
	}

	err := s.repo.User.Create(user)
	if err != nil {
		s.log.Error("Error creating user", zap.Error(err))
		return err
	}

	if user.Role == "admin" {
		adminPermission := domain.UserPermission{
			UserID:       user.ID,
			PermissionID: 1,
		}

		err = s.repo.UserPermission.Create(adminPermission)
		if err != nil {
			s.log.Error("Error creating admin permission", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *userService) Get(user domain.User) (*domain.User, error) {
	return s.repo.User.Get(user)
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

func (s *userService) Delete(id uint) error {
	return s.repo.User.Delete(id)
}

func (s *userService) Update(user domain.User) error {
	existedUser := s.repo.User.GetByEmail(user.Email)
	if existedUser != nil && existedUser.ID != user.ID {
		return errors.New("user email already exists")
	}

	return s.repo.User.Update(&user)
}
