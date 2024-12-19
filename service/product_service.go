package service

import (
	"errors"
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type ProductService interface {
	All(page, limit int, productStatus, categoryName, stock string, quantity int, minPrice, maxPrice float64) ([]*domain.Product, int64, error)
	Add(input *domain.Product, categoryName string) (*domain.Product, error)
	Update(id uint, ProductData *domain.Product, categoryName string) (*domain.Product, error)
	Delete(id uint) error
}

type productService struct {
	repo repository.ProductRepository
	log  *zap.Logger
}

func NewProductService(repo repository.ProductRepository, log *zap.Logger) ProductService {
	return &productService{repo, log}
}

func (s *productService) All(page, limit int, productStatus, categoryName, stock string, quantity int, minPrice, maxPrice float64) ([]*domain.Product, int64, error) {

	// Validasi awal untuk filter jika diperlukan
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Panggil repository untuk mendapatkan data inventori
	product, totalItems, err := s.repo.All(page, limit, productStatus, categoryName, stock, quantity, minPrice, maxPrice)
	if err != nil {
		s.log.Error("Failed to fetch product", zap.Error(err))
		return nil, 0, err
	}

	if len(product) == 0 {
		s.log.Warn("No product found", zap.Int("page", page), zap.Int("limit", limit))
		return nil, 0, errors.New("product not found")
	}

	return product, totalItems, nil
}

func (s *productService) Add(input *domain.Product, categoryName string) (*domain.Product, error) {
	// Validasi awal input
	if categoryName == "" {
		return nil, errors.New("category name cannot be empty")
	}

	if input.Stock <= 0 {
		s.log.Warn("Invalid Stock", zap.Int("Stock", input.Stock))
		return nil, errors.New("stock must be greater than 0")
	}
	if input.Price <= 0 {
		s.log.Warn("Invalid Price", zap.Float64("price", input.Price))
		return nil, errors.New("price must be greater than 0")
	}

	// Panggil repository untuk menyimpan data
	result, err := s.repo.Add(input, categoryName)
	if err != nil {
		s.log.Error("Failed to add inventory", zap.Error(err))
		return nil, err
	}

	s.log.Info("Inventory added successfully", zap.Uint("inventory_id", result.ID))
	return result, nil
}

func (s *productService) Update(id uint, inventoryData *domain.Product, categoryName string) (*domain.Product, error) {

	// Panggil repository untuk update inventory
	updatedInventory, err := s.repo.Update(id, inventoryData, categoryName)
	if err != nil {
		s.log.Error("Failed to update inventory", zap.Error(err))
		return nil, err
	}

	return updatedInventory, nil
}

func (s *productService) Delete(id uint) error {
	// Panggil repository untuk soft delete inventory
	err := s.repo.Delete(id)
	if err != nil {
		s.log.Error("Failed to soft delete inventory", zap.Error(err))
		return err
	}

	return nil
}
