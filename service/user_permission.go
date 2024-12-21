package service

import (
	"go.uber.org/zap"
	"project/domain"
	"project/repository"
)

type UserPermissionService interface {
	Update(userID uint, permissions []uint) error
}
type userPermissionService struct {
	repo repository.UserPermissionRepository
	log  *zap.Logger
}

func NewUserPermissionService(repo repository.UserPermissionRepository, log *zap.Logger) UserPermissionService {
	return &userPermissionService{repo, log}
}

func (s *userPermissionService) Update(userID uint, permissions []uint) error {
	var newPermissions []domain.Permission
	for _, permission := range permissions {
		newPermissions = append(newPermissions, domain.Permission{ID: permission})
	}
	user := domain.User{ID: userID, Permissions: newPermissions}
	return s.repo.Update(user)
}
