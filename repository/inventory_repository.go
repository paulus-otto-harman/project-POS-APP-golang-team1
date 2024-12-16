package repository

import (
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InventoryRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewInventoryRepository(db *gorm.DB, log *zap.Logger) *InventoryRepository {
	return &InventoryRepository{db: db, log: log}
}

func (repo InventoryRepository) All(page, limit int, categoryID string) ([]*domain.Inventory, int64, error) {
	var Inventory []*domain.Inventory
	var totalItems int64

	query := repo.db.Model(&domain.Inventory{})
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		repo.log.Error("Failed to count total inventory", zap.Error(err))
		return nil, 0, err
	}

	if totalItems == 0 {
		repo.log.Warn("No inventory found")
		return []*domain.Inventory{}, 0, nil
	}

	err := query.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, created_at, updated_at")
	}).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&Inventory).Error
	if err != nil {
		repo.log.Error("Failed to fetch Inventory", zap.Error(err))
		return nil, 0, err
	}

	return Inventory, totalItems, nil
}
