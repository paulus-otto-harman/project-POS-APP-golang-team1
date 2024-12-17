package service

import (
	"errors"
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type InventoryService interface {
	All(page, limit int, productStatus, categoryName, stock string, quantity int, minPrice, maxPrice float64) ([]*domain.Inventory, int64, error)
	Add(input *domain.Inventory, categoryName string) (*domain.Inventory, error)
	Update(id uint, inventoryData *domain.Inventory, categoryName string) (*domain.Inventory, error)
}

type inventoryService struct {
	repo repository.InventoryRepository
	log  *zap.Logger
}

func NewInventoryService(repo repository.InventoryRepository, log *zap.Logger) InventoryService {
	return &inventoryService{repo, log}
}

func (s *inventoryService) All(page, limit int, productStatus, categoryName, stock string, quantity int, minPrice, maxPrice float64) ([]*domain.Inventory, int64, error) {

	// Validasi awal untuk filter jika diperlukan
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Panggil repository untuk mendapatkan data inventori
	inventory, totalItems, err := s.repo.All(page, limit, productStatus, categoryName, stock, quantity, minPrice, maxPrice)
	if err != nil {
		s.log.Error("Failed to fetch inventory", zap.Error(err))
		return nil, 0, err
	}

	if len(inventory) == 0 {
		s.log.Warn("No inventory found", zap.Int("page", page), zap.Int("limit", limit))
		return nil, 0, errors.New("inventory not found")
	}

	return inventory, totalItems, nil
}

func (s *inventoryService) Add(input *domain.Inventory, categoryName string) (*domain.Inventory, error) {
	// Validasi awal input
	if categoryName == "" {
		return nil, errors.New("category name cannot be empty")
	}

	if input.Quantity <= 0 {
		s.log.Warn("Invalid Quantity", zap.Int("quantity", input.Quantity))
		return nil, errors.New("quantity must be greater than 0")
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

func (s *inventoryService) Update(id uint, inventoryData *domain.Inventory, categoryName string) (*domain.Inventory, error) {

	// Panggil repository untuk update inventory
	updatedInventory, err := s.repo.Update(id, inventoryData, categoryName)
	if err != nil {
		s.log.Error("Failed to update inventory", zap.Error(err))
		return nil, err
	}

	return updatedInventory, nil
}
