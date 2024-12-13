package service

import (
	"errors"
	"io"
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type CategoryService interface {
	All(page, limit int) ([]*domain.Category, int64, error)
	Create(category *domain.Category, file io.Reader, filename string) error
}

type categoryService struct {
	repo repository.CategoryRepository
	log  *zap.Logger
}

func NewCategoryService(repo repository.CategoryRepository, log *zap.Logger) CategoryService {
	return &categoryService{repo, log}
}

func (s *categoryService) All(page, limit int) ([]*domain.Category, int64, error) {
	categories, totalItems, err := s.repo.All(page, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(categories) == 0 {
		return nil, int64(totalItems), errors.New("categories not found")
	}

	return categories, int64(totalItems), nil
}
func (s *categoryService) Create(category *domain.Category, file io.Reader, filename string) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}

	return s.repo.Create(category, file, filename)
}
