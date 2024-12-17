package repository

import (
	"errors"
	// "fmt"
	// "log"

	// "math"
	"project/domain"
	// "project/helper"

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

// func (repo OrderRepository) GenerateCodeOrder() (string, error) {
// 	var lastOrder domain.Order
// 	if err := repo.db.Order("id desc").First(&lastOrder).Error; err != nil && err != gorm.ErrRecordNotFound {
// 		return "", errors.New("failed to retrieve last order")
// 	}

// 	lastCode := "ORD000"
// 	if lastOrder.CodeOrder != "" {
// 		lastCode = lastOrder.CodeOrder
// 		log.Println(lastCode, "<<<")
// 	}
	
// 	codeNum, err := helper.Uint(lastCode[3:])
// 	if err != nil {
// 		return "", errors.New("failed to parse last order code number")
// 	}
// 	log.Println(codeNum, ">>>>>")
	
// 	newCodeNum := codeNum + 1
// 	newCode := fmt.Sprintf("ORD%03d", newCodeNum)
// 	log.Println(newCode, ">>>>>")

// 	return newCode, nil
// }

func (repo OrderRepository) Create(Order *domain.Order) error {

	err := repo.db.Create(&Order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			repo.log.Error("Duplicate code order", zap.Error(err))
			return errors.New("order with this code order already exists")
		}
		repo.log.Error("Failed to save Order", zap.Error(err))
		return err
	}
	// Order.CodeOrder, err = repo.GenerateCodeOrder()
	// if err != nil {
	// 	return err
	// }

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

func (repo *OrderRepository) FindByIDOrder(order *domain.OrderDetail, id uint) error {
	if err := repo.db.First(order, "order_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		repo.log.Error("Failed to fetch Order by ID", zap.Error(err))
		return err
	}
	return nil
}
func (repo *OrderRepository) FindByIDTable(table *domain.Table, id string) error {
	if err := repo.db.First(table, "id = ?", id).Error; err != nil {
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

func (repo OrderRepository) AllOrders(page, limit int, name, codeOrder, status string) ([]*domain.OrderDetail, int64, error) {
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
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&totalItems).Error; err != nil {
		repo.log.Error("Failed to count total orders", zap.Error(err))
		return nil, 0, err
	}

	if totalItems == 0 {
		repo.log.Warn("No orders found")
		return nil, 0, nil
	}

	err := query.
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&orders).Error
	if err != nil {
		repo.log.Error("Failed to fetch orders", zap.Error(err))
		return nil, 0, err
	}

	return orders, totalItems, nil
}
