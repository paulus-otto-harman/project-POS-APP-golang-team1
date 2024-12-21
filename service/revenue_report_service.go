package service

import (
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type RevenueService interface {
	GetTotalRevenueByStatus() (map[string]interface{}, error)
	GetMonthlyRevenue() (map[string]float64, error)
	GetProductRevenueDetails() ([]*domain.ProductRevenue, error)
}

type revenueService struct {
	repo repository.RevenueRepository
	log  *zap.Logger
}

func NewRevenueService(repo repository.RevenueRepository, log  *zap.Logger) RevenueService {
	return &revenueService{repo: repo, log: log}
}

func (s *revenueService) GetTotalRevenueByStatus() (map[string]interface{}, error) {
	return s.repo.GetTotalRevenueByStatus()
}

func (s *revenueService) GetMonthlyRevenue() (map[string]float64, error) {
	return s.repo.GetMonthlyRevenue()
}

func (s *revenueService) GetProductRevenueDetails() ([]*domain.ProductRevenue, error) {
	return s.repo.GetProductRevenueDetails()
}
