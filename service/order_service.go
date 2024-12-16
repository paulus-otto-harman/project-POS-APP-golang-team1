package service

import (
	"errors"
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type OrderService interface {
	AllTables(page, limit int) ([]*domain.Table, int64, error)
	Create(order *domain.Order) error
	FindByID(order *domain.Order, id string) error
	Update(order *domain.Order) error
	AllOrders(page, limit int, orderID string) ([]*domain.Order, int64, error)
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
func (s *orderService) Create(order *domain.Order) error {
	if order.Name == "" {
		return errors.New("order name is required")
	}

	return s.repo.Create(order)
}

func (s *orderService) FindByID(order *domain.Order, id string) error {
	if err := s.repo.FindByID(order, id); err != nil {
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
func (s *orderService) AllOrders(page, limit int, orderID string) ([]*domain.Order, int64, error) {

	orders, totalItems, err := s.repo.AllOrders(page, limit, orderID)
	if err != nil {
		return nil, 0, err
	}
	if len(orders) == 0 {
		return nil, int64(totalItems), errors.New("orders not found")
	}

	return orders, int64(totalItems), nil
}
