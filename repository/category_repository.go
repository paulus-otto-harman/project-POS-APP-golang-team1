package repository

import (
	"errors"
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewCategoryRepository(db *gorm.DB, log *zap.Logger) *CategoryRepository {
	return &CategoryRepository{db: db, log: log}
}

func (repo CategoryRepository) Create(category *domain.Category) error {

	err := repo.db.Create(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			repo.log.Error("Duplicate category name", zap.Error(err))
			return errors.New("category with this name already exists")
		}
		repo.log.Error("Failed to save category", zap.Error(err))
		return err
	}

	repo.log.Info("Category successfully created")
	return nil
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

func (repo *CategoryRepository) FindByID(category *domain.Category, id string) error {
	if err := repo.db.First(category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		repo.log.Error("Failed to fetch category by ID", zap.Error(err))
		return err
	}
	return nil
}

func (repo *CategoryRepository) Update(category *domain.Category) error {
	if err := repo.db.Save(category).Error; err != nil {
		repo.log.Error("Failed to update category", zap.Error(err))
		return err
	}
	return nil
}

func (repo CategoryRepository) AllProducts(page, limit int, categoryID string) ([]*domain.Inventory, int64, error) {
	var products []*domain.Inventory
	var totalItems int64

	query := repo.db.Model(&domain.Inventory{})
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		repo.log.Error("Failed to count total products", zap.Error(err))
		return nil, 0, err
	}

	if totalItems == 0 {
		repo.log.Warn("No products found")
		return []*domain.Inventory{}, 0, nil
	}

	err := query.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, created_at, updated_at")
	}).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&products).Error
	if err != nil {
		repo.log.Error("Failed to fetch products", zap.Error(err))
		return nil, 0, err
	}

	return products, totalItems, nil
}
