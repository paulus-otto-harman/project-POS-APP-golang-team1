package repository

import (
	"errors"
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewOrderRepository(db *gorm.DB, log *zap.Logger) *OrderRepository {
	return &OrderRepository{db: db, log: log}
}

func (repo OrderRepository) Create(Order *domain.Order) error {

	err := repo.db.Create(&Order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			repo.log.Error("Duplicate Order name", zap.Error(err))
			return errors.New("order with this name already exists")
		}
		repo.log.Error("Failed to save Order", zap.Error(err))
		return err
	}

	repo.log.Info("Order successfully created")
	return nil
}

func (repo OrderRepository) AllTables(page, limit int) ([]*domain.Table, int64, error) {
	var tables []*domain.Table
	var totalItems int64

	offset := (page - 1) * limit

	err := repo.db.Model(&domain.Table{}).Where("status = ?", true).Count(&totalItems).Error
	if err != nil {
		repo.log.Error("Failed to count total tables", zap.Error(err))
		return nil, 0, err
	}

	err = repo.db.Model(&domain.Table{}).Where("status = ?", true).
		Offset(offset).
		Limit(limit).
		Find(&tables).Error
	if err != nil {
		repo.log.Error("Failed to fetch tables", zap.Error(err))
		return nil, 0, err
	}

	if totalItems == 0 {
		repo.log.Warn("No tables found")
		return nil, 0, errors.New("no tables found")
	}

	return tables, totalItems, nil
}

func (repo *OrderRepository) FindByID(Order *domain.Order, id string) error {
	if err := repo.db.First(Order, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		repo.log.Error("Failed to fetch Order by ID", zap.Error(err))
		return err
	}
	return nil
}

func (repo *OrderRepository) Update(Order *domain.Order) error {
	if err := repo.db.Save(Order).Error; err != nil {
		repo.log.Error("Failed to update Order", zap.Error(err))
		return err
	}
	return nil
}

func (repo OrderRepository) AllOrders(page, limit int, OrderID string) ([]*domain.Order, int64, error) {
	var Orders []*domain.Order
	var totalItems int64

	query := repo.db.Model(&domain.Order{})
	if OrderID != "" {
		query = query.Where("Order_id = ?", OrderID)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		repo.log.Error("Failed to count total Orders", zap.Error(err))
		return nil, 0, err
	}

	if totalItems == 0 {
		repo.log.Warn("No Orders found")
		return []*domain.Order{}, 0, nil
	}

	err := query.Preload("Order", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, created_at, updated_at")
	}).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&Orders).Error
	if err != nil {
		repo.log.Error("Failed to fetch Orders", zap.Error(err))
		return nil, 0, err
	}

	return Orders, totalItems, nil
}
