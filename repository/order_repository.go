package repository

import (
	"errors"
	"project/domain"
	"project/helper"

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

func (repo OrderRepository) Create(order *domain.Order) error {

	return repo.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&order).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				repo.log.Error("Duplicate code order", zap.Error(err))
				return errors.New("order with this code order already exists")
			}
			repo.log.Error("Failed to save Order", zap.Error(err))
			return err
		}

		repo.log.Info("Order successfully created")
		return nil
	})
}

func (repo OrderRepository) AllTables(page, limit int) ([]*domain.Table, int64, error) {
	var tables []*domain.Table
	var totalItems int64

	err := repo.db.Model(&domain.Table{}).Where("status = ?", true).Count(&totalItems).Error
	if err != nil {
		repo.log.Error("Failed to count total tables", zap.Error(err))
		return nil, 0, err
	}

	err = repo.db.Model(&domain.Table{}).Order("id").Where("status = ?", true).
		Scopes(helper.Paginate(uint(page), uint(limit))).
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
func (repo OrderRepository) AllPayments() ([]*domain.PaymentMethod, error) {
	var payments []*domain.PaymentMethod

	err := repo.db.Model(&domain.PaymentMethod{}).Where("status = ?", true).
		Find(&payments).Error
	if err != nil {
		repo.log.Error("Failed to fetch payments", zap.Error(err))
		return nil, err
	}

	if len(payments) == 0 {
		repo.log.Warn("No payments found")
		return nil, errors.New("no payments found")
	}

	return payments, nil
}

func (repo *OrderRepository) FindByIDOrder(order *domain.Order, id string) error {
	if err := repo.db.First(order, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		repo.log.Error("Failed to fetch Order by ID", zap.Error(err))
		return err
	}
	return nil
}
func (repo *OrderRepository) FindByIDOrderDetail(order *domain.OrderDetail, id string) error {
	if err := repo.db.First(order, "order_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		repo.log.Error("Failed to fetch Order by ID", zap.Error(err))
		return err
	}
	return nil
}

// func (repo *OrderRepository) FindByIDTable(table *domain.Table, id string) error {
// 	if err := repo.db.First(table, "id = ?", id).Where("status", true).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return errors.New("table not found")
// 		}
// 		repo.log.Error("Failed to fetch table by ID", zap.Error(err))
// 		return err
// 	}
// 	return nil
// }

func (repo *OrderRepository) Update(order *domain.Order) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		
		// tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&order)
		if err := tx.Model(&order).Where("id = ?", order.ID).Updates(order).Error; err != nil {
			repo.log.Error("Failed to update Order", zap.Error(err))
			return err
		}
		return nil
	})
}

func (repo OrderRepository) AllOrders(page, limit int, name, codeOrder string, status domain.StatusPayment) ([]*domain.OrderDetail, int64, error) {
	var orders []*domain.OrderDetail
	var totalItems int64

	query := repo.db.Model(&domain.OrderDetail{})
	if name != "" {
		query = query.Where("name ILIKE ?", name+"%")
	}
	if codeOrder != "" {
		query = query.Where("code_order = ?", codeOrder)
	}
	if status != "" {
		query = query.Where("status_payment = ?", status)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		repo.log.Error("Failed to count total orders", zap.Error(err))
		return nil, 0, err
	}

	if totalItems == 0 {
		repo.log.Warn("No orders found")
		return nil, 0, nil
	}

	err := query.Order("order_id").Scopes(helper.Paginate(uint(page), uint(limit))).
		Find(&orders).Error
	if err != nil {
		repo.log.Error("Failed to fetch orders", zap.Error(err))
		return nil, 0, err
	}

	return orders, totalItems, nil
}

func (repo *OrderRepository) Delete(order *domain.Order) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&order, order.ID).Error; err != nil {
			repo.log.Error("Error deleting notification", zap.Error(err))
			return err
		}
		return nil
	})
}
