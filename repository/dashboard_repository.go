package repository

import (
	"fmt"
	"project/domain"
	"time"

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
	var totalTables int64

	// Current date and month
	// For correct timezone handling, you can use time.Now().UTC() for UTC time
	today := time.Now().Format("2006-01-02") // today without time zone
	month := time.Now().Format("2006-01")    // current month for monthly sales query

	// Log today and month
	repo.log.Info("Today in Jakarta Timezone", zap.String("today", today))
	repo.log.Info("Month in Jakarta Timezone", zap.String("month", month))

	// Daily Sales Query
	repo.log.Info("Executing Daily Sales Query", zap.String("query", today))
	err := repo.db.Model(&domain.Order{}).
		Where("DATE(orders.created_at) = ? AND orders.status_payment = ?", today, domain.OrderCompleted).
		Select("COALESCE(SUM(order_items.quantity * products.price), 0) AS daily_sales").
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Scan(&dailySales).Error

	if err != nil {
		repo.log.Error("Error executing daily sales query", zap.Error(err))
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
		repo.log.Error("Error executing monthly sales query", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch monthly sales: %v", err)
	}

	// Total Tables Query
	err = repo.db.Model(&domain.Table{}).Count(&totalTables).Error
	if err != nil {
		repo.log.Error("Error fetching total tables", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch total tables: %v", err)
	}

	// Setting the data to the summary
	summary.DailySales = dailySales
	summary.MonthlySales = monthlySales
	summary.TotalTables = int(totalTables)

	return &summary, nil
}
