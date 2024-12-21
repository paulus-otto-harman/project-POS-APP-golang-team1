package repository

import (
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RevenueRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewRevenueRepository(db *gorm.DB, log *zap.Logger) *RevenueRepository {
	return &RevenueRepository{db: db, log: log}
}

func (repo RevenueRepository) GetTotalRevenueByStatus() (map[string]float64, error) {
	result := make(map[string]float64)

	return result, nil
}

func (repo RevenueRepository) GetMonthlyRevenue() (map[string]float64, error) {
	result := make(map[string]float64)

	return result, nil
}

func (repo RevenueRepository) GetProductRevenueDetails() ([]*domain.ProductRevenue, error) {
	var products []*domain.ProductRevenue

	return products, nil
}
