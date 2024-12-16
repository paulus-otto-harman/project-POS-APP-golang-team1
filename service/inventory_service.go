package service

import (
	"errors"
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type InventoryService interface {
	All(page, limit int, categoryID string) ([]*domain.Inventory, int64, error)
}

type inventoryService struct {
	repo repository.InventoryRepository
	log  *zap.Logger
}

func NewInventoryService(repo repository.InventoryRepository, log *zap.Logger) InventoryService {
	return &inventoryService{repo, log}
}

func (s *inventoryService) All(page, limit int, categoryID string) ([]*domain.Inventory, int64, error) {

	products, totalItems, err := s.repo.All(page, limit, categoryID)
	if err != nil {
		return nil, 0, err
	}
	if len(products) == 0 {
		return nil, int64(totalItems), errors.New("inventory not found")
	}

	return products, int64(totalItems), nil
}
