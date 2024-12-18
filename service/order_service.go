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
	FindByIDOrder(order *domain.OrderDetail, id uint) error
	Update(order *domain.Order) error
	AllOrders(page, limit int, name, codeOrder, status string) ([]*domain.OrderDetail, int64, error)
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

func (s *orderService) FindByIDOrder(order *domain.OrderDetail, id uint) error {
	if err := s.repo.FindByIDOrder(order, id); err != nil {
		s.log.Error("Failed to find order", zap.Error(err))
		return err
	}
	return nil
}

func (s *orderService) Update(order *domain.Order) error {
	if order.Name == "" {
		return errors.New("order name is required")
	}

	if err := s.repo.Update(order); err != nil {
		s.log.Error("Failed to update order", zap.Error(err))
		return err
	}
	return nil
}
func (s *orderService) AllOrders(page, limit int, name, codeOrder, status string) ([]*domain.OrderDetail, int64, error) {
	orders, totalItems, err := s.repo.AllOrders(page, limit, name, codeOrder, status)
	if err != nil {
		return nil, 0, err
	}
	if len(orders) == 0 {
		return nil, int64(totalItems), errors.New("orders not found")
	}

	return orders, int64(totalItems), nil
}
