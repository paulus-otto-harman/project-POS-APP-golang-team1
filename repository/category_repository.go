package repository

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project/domain"
)

type CategoryRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewCategoryRepository(db *gorm.DB, log *zap.Logger) *CategoryRepository {
	return &CategoryRepository{db: db, log: log}
}

func (repo CategoryRepository) Create(category *domain.Category) error {
	return repo.db.Create(&category).Error
}

func (repo CategoryRepository) All(page, limit int) ([]*domain.Category, int64, error) {
	var categories []*domain.Category
	var totalItems int64

	offset := (page - 1) * limit

	err := repo.db.Model(&domain.Category{}).Count(&totalItems).Error
	if err != nil {
		repo.log.Error("Failed to count total categories", zap.Error(err))
		return nil, 0, err
	}

	err = repo.db.Model(&domain.Category{}).
		Offset(offset).
		Limit(limit).
		Find(&categories).Error
	if err != nil {
		repo.log.Error("Failed to fetch categories", zap.Error(err))
		return nil, 0, err
	}

	if len(categories) == 0 {
		repo.log.Warn("No categories found")
		return nil, 0, errors.New("no categories found")
	}

	return categories, totalItems, nil
}
