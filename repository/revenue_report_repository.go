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

func (repo RevenueRepository) GetTotalRevenueByStatus() (map[string]interface{}, error) {
	var results []struct {
        StatusPayment string  `json:"status_payment"`
        Revenue       float64 `json:"revenue"`
    }

    err := repo.db.Model(&domain.OrderDetail{}).
        Select("status_payment, SUM(total) as revenue").
        Group("status_payment").
        Scan(&results).Error
    if err != nil {
        return nil, err
    }

	var totalRevenue float64
    revenueMap := make(map[string]float64)
    for _, result := range results {
		totalRevenue += result.Revenue
        revenueMap[result.StatusPayment] = result.Revenue

    }

	response := map[string]interface{}{
        "total_revenue": totalRevenue,
        "by_status":     revenueMap,
    }

    return response, nil
}

func (repo RevenueRepository) GetMonthlyRevenue() (map[string]float64, error) {
	result := make(map[string]float64)

	return result, nil
}

func (repo RevenueRepository) GetProductRevenueDetails() ([]*domain.ProductRevenue, error) {
	var products []*domain.ProductRevenue

	return products, nil
}
