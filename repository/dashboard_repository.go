package repository

import (
	"fmt"
	"time"

	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DashboardRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewDashboardRepository(db *gorm.DB, log *zap.Logger) *DashboardRepository {
	return &DashboardRepository{db: db, log: log}
}

func (repo *DashboardRepository) GetDashboardSummary() (*domain.Dashboard, error) {
	var summary domain.Dashboard
	var dailySales, monthlySales float64
	var totalTables, occupiedTables int64

	// Current date and month
	today := time.Now().UTC().Format("2006-01-02")
	month := time.Now().Format("2006-01")

	// Daily Sales Query
	err := repo.db.Model(&domain.Order{}).
		Where("DATE(orders.created_at) = ? AND orders.status_payment = ?", today, domain.OrderCompleted).
		Select("COALESCE(SUM(order_items.quantity * products.price), 0) AS daily_sales").
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Scan(&dailySales).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily sales: %v", err)
	}

	// Monthly Sales Query
	err = repo.db.Model(&domain.Order{}).
		Where("TO_CHAR(orders.created_at, 'YYYY-MM') = ? AND orders.status_payment = ?", month, domain.OrderCompleted).
		Select("COALESCE(SUM(order_items.quantity * products.price), 0) AS monthly_sales").
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Scan(&monthlySales).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch monthly sales: %v", err)
	}

	// Total Tables Query
	err = repo.db.Model(&domain.Table{}).Count(&totalTables).Error
	if err != nil {
		repo.log.Error("Error fetching total tables", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch total tables: %v", err)
	}

	// All Time Table Occupancy Query
	err = repo.db.Model(&domain.Order{}).
		Joins("JOIN tables ON orders.table_id = tables.id").
		Where("orders.status_payment = ?", domain.OrderCompleted).
		Count(&occupiedTables).Error
	if err != nil {
		repo.log.Error("Error fetching occupied tables", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch occupied tables: %v", err)
	}

	// Calculate occupancy percentage
	occupancyPercentage := (float64(occupiedTables) / float64(totalTables)) * 100

	// Set summary
	summary.DailySales = dailySales
	summary.MonthlySales = monthlySales
	summary.TableOccupancy = occupancyPercentage

	return &summary, nil
}
