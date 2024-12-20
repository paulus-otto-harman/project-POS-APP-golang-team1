package service

import (
	"go.uber.org/zap"
	"project/repository"
)

type UserPermissionService interface {
	Update() error
}
type userPermissionService struct {
	repo repository.UserPermissionRepository
	log  *zap.Logger
}

func NewUserPermissionService(repo repository.UserPermissionRepository, log *zap.Logger) UserPermissionService {
	return &userPermissionService{repo, log}
}

func (s *userPermissionService) Update() error {
	return nil
}
