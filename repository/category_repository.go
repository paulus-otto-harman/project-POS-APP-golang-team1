package repository

import (
	"errors"
	"io"
	"project/domain"
	"project/helper"

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

func (repo CategoryRepository) Create(category *domain.Category, file io.Reader, filename string) error {

	imageURL, err := helper.UploadFileThirdPartyAPI(file, filename)
	if err != nil {
		repo.log.Error("Failed to upload file", zap.Error(err))
		return err
	}

	category.Icon = imageURL

	err = repo.db.Create(&category).Error
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
