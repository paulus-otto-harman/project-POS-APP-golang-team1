package service

import (
	"go.uber.org/zap"
	"project/domain"
	"project/repository"
)

type UserService interface {
	All(user domain.User) ([]domain.User, error)
	Get(user domain.User) (*domain.User, error)
	Register(user *domain.User) error
}

type userService struct {
	repo repository.UserRepository
	log  *zap.Logger
}

func NewUserService(repo repository.UserRepository, log *zap.Logger) UserService {
	return &userService{repo, log}
}

func (s *userService) All(user domain.User) ([]domain.User, error) {
	return s.repo.All(user)
}

func (s *userService) Get(user domain.User) (*domain.User, error) {
	return s.repo.Get(user)
}

func (s *userService) Register(user *domain.User) error {
	return s.repo.Create(user)
}
