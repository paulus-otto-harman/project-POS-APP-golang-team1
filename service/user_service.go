package service

import (
	"errors"
	"project/domain"
	"project/helper"
	"project/repository"

	"go.uber.org/zap"
)

type UserService interface {
	All(sortField, sortDirection string, page, limit uint) ([]domain.User, int64, error)
	Register(user *domain.User) error
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
