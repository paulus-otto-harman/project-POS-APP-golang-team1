package service

import (
	"errors"
	"io"
	"project/domain"
	"project/helper"
	"project/repository"

	"go.uber.org/zap"
)

type CategoryService interface {
	All(page, limit int) ([]*domain.Category, int64, error)
	Create(category *domain.Category) error
	FindByID(category *domain.Category, id string) error
	Update(category *domain.Category) error
	UploadIcon(file io.Reader, filename string) (string, error)
	AllProducts(page, limit int, categoryID string) ([]*domain.Product, int64, error)
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
func (s *categoryService) Create(category *domain.Category) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}

	return s.repo.Create(category)
}

func (s *categoryService) FindByID(category *domain.Category, id string) error {
	if err := s.repo.FindByID(category, id); err != nil {
		s.log.Error("Failed to find category", zap.Error(err))
		return err
	}
	return nil
}

func (s *categoryService) UploadIcon(file io.Reader, filename string) (string, error) {
	return helper.UploadFileThirdPartyAPI(file, filename)
}

func (s *categoryService) Update(category *domain.Category) error {
	if category.Name == "" {
		return errors.New("category name is required")
	}

	if err := s.repo.Update(category); err != nil {
		s.log.Error("Failed to update category", zap.Error(err))
		return err
	}
	return nil
}
func (s *categoryService) AllProducts(page, limit int, categoryID string) ([]*domain.Product, int64, error) {

	products, totalItems, err := s.repo.AllProducts(page, limit, categoryID)
	if err != nil {
		return nil, 0, err
	}
	if len(products) == 0 {
		return nil, int64(totalItems), errors.New("products not found")
	}

	return products, int64(totalItems), nil
}
