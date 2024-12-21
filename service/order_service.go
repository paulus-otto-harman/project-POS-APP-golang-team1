package service

import (
	"errors"
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type OrderService interface {
	AllTables(page, limit int) ([]*domain.Table, int64, error)
	AllPayments() ([]*domain.PaymentMethod, error)
	CreateOrder(name string, tableID uint, orderItems []domain.OrderItem) error
	FindByIDOrder(order *domain.Order, id string) error
	// FindByIDTable(table *domain.Table, id string) error
	FindByIDOrderDetail(order *domain.OrderDetail, id string) error
	Update(order *domain.Order) error
	AllOrders(page, limit int, name, codeOrder string, status domain.StatusPayment) ([]*domain.OrderDetail, int64, error)
	Delete(order *domain.Order) error
}

type orderService struct {
	repo repository.OrderRepository
	log  *zap.Logger
}

func NewOrderService(repo repository.OrderRepository, log *zap.Logger) OrderService {
	return &orderService{repo, log}
}

func (s *orderService) AllTables(page, limit int) ([]*domain.Table, int64, error) {
	tables, totalItems, err := s.repo.AllTables(page, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(tables) == 0 {
		return nil, int64(totalItems), errors.New("tables not found")
	}

	return tables, int64(totalItems), nil
}
func (s *orderService) AllPayments() ([]*domain.PaymentMethod, error) {
	payments, err := s.repo.AllPayments()
	if len(payments) == 0 {
		return nil, errors.New("payments not found")
	}
	if err != nil {
		return nil, err
	}

	return payments, nil
}
func (s *orderService) CreateOrder(name string, tableID uint, orderItems []domain.OrderItem) error {
	if len(orderItems) == 0 {
		return errors.New("order items cannot be empty")
	}
	order := &domain.Order{
		Name:       name,
		TableID:    tableID,
		OrderItems: orderItems,
	}

	if err := s.repo.Create(order); err != nil {
		return err
	}

	return nil
}

func (s *orderService) FindByIDOrder(order *domain.Order, id string) error {
	if err := s.repo.FindByIDOrder(order, id); err != nil {
		s.log.Error("Failed to find order", zap.Error(err))
		return err
	}
	return nil
}
func (s *orderService) FindByIDOrderDetail(order *domain.OrderDetail, id string) error {
	if err := s.repo.FindByIDOrderDetail(order, id); err != nil {
		s.log.Error("Failed to find order", zap.Error(err))
		return err
	}
	return nil
}
// func (s *orderService) FindByIDTable(table *domain.Table, id string) error {
// 	if err := s.repo.FindByIDTable(table, id); err != nil {
// 		s.log.Error("Failed to find order", zap.Error(err))
// 		return err
// 	}
// 	return nil
// }

func (s *orderService) Update(order *domain.Order) error {

	if len(order.OrderItems) == 0 {
		return errors.New("order items cannot be empty")
	}

	if err := s.repo.Update(order); err != nil {
		s.log.Error("Failed to update order", zap.Error(err))
		return err
	}
	return nil
}
func (s *orderService) AllOrders(page, limit int, name, codeOrder string, status domain.StatusPayment) ([]*domain.OrderDetail, int64, error) {
	orders, totalItems, err := s.repo.AllOrders(page, limit, name, codeOrder, status)
	if err != nil {
		return nil, 0, err
	}
	if len(orders) == 0 {
		return nil, int64(totalItems), errors.New("orders not found")
	}

	return orders, int64(totalItems), nil
}

func (s *orderService) Delete(order *domain.Order) error {
	if err := s.repo.Delete(order); err != nil {
		s.log.Error("Failed to delete order", zap.Error(err))
		return err
	}
	return nil
}
