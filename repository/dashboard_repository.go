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

// FetchSummary retrieves summary data (daily sales, monthly sales, and table occupancy).
func (repo *DashboardRepository) FetchSummary() (float64, float64, float64, error) {
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
		return 0, 0, 0, fmt.Errorf("failed to fetch daily sales: %v", err)
	}

	// Monthly Sales Query
	err = repo.db.Model(&domain.Order{}).
		Where("TO_CHAR(orders.created_at, 'YYYY-MM') = ? AND orders.status_payment = ?", month, domain.OrderCompleted).
		Select("COALESCE(SUM(order_items.quantity * products.price), 0) AS monthly_sales").
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Scan(&monthlySales).Error
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to fetch monthly sales: %v", err)
	}

	// Total Tables Query
	err = repo.db.Model(&domain.Table{}).Count(&totalTables).Error
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to fetch total tables: %v", err)
	}

	// Occupied Tables Query
	err = repo.db.Model(&domain.Order{}).
		Where("orders.status_payment = ?", domain.OrderCompleted).
		Joins("JOIN tables ON orders.table_id = tables.id").
		Count(&occupiedTables).Error
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to fetch occupied tables: %v", err)
	}

	// Calculate occupancy percentage
	occupancyPercentage := (float64(occupiedTables) / float64(totalTables)) * 100

	return dailySales, monthlySales, occupancyPercentage, nil
}

// FetchPopularProducts retrieves the most popular products.
func (repo *DashboardRepository) FetchPopularProducts() ([]domain.PopularNewResponse, error) {
	var popularProducts []domain.PopularNewResponse

	err := repo.db.Model(&domain.Order{}).
		Where("orders.status_payment = ?", domain.OrderCompleted).
		Select(`
			products.name,
			COUNT(order_items.product_id) AS order_count,
			products.image,
			products.availability,
			products.price
		`).
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Joins("JOIN products ON order_items.product_id = products.id").
		Group("products.id, products.name, products.image, products.availability, products.price").
		Order("order_count DESC").
		Limit(4).
		Scan(&popularProducts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch popular products: %v", err)
	}

	return popularProducts, nil
}

// FetchNewProducts retrieves new products added within the last 30 days.
func (repo *DashboardRepository) FetchNewProducts() ([]domain.PopularNewResponse, error) {
	var newProducts []domain.PopularNewResponse
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	err := repo.db.Model(&domain.Product{}).
		Where("products.created_at >= ?", thirtyDaysAgo).
		Select(`
			products.name,
			0 AS order_count,
			products.image,
			products.availability,
			products.price
		`).
		Order("products.created_at DESC").
		Limit(4).
		Scan(&newProducts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch new products: %v", err)
	}

	return newProducts, nil
}

// GetDashboard aggregates data for the dashboard.
func (repo *DashboardRepository) GetDashboard() (*domain.Dashboard, error) {
	// Fetch summary
	dailySales, monthlySales, tableOccupancy, err := repo.FetchSummary()
	if err != nil {
		return nil, err
	}

	// Fetch popular products
	popularProducts, err := repo.FetchPopularProducts()
	if err != nil {
		return nil, err
	}

	// Fetch new products
	newProducts, err := repo.FetchNewProducts()
	if err != nil {
		return nil, err
	}

	// Count monthly orders
	month := time.Now().Format("2006-01")
	var monthlyOrderCount int64
	err = repo.db.Model(&domain.Order{}).
		Where("TO_CHAR(orders.created_at, 'YYYY-MM') = ? AND orders.status_payment = ?", month, domain.OrderCompleted).
		Count(&monthlyOrderCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch monthly order count: %v", err)
	}

	// Set summary
	summary := &domain.Dashboard{
		DailySales:        dailySales,
		MonthlySales:      monthlySales,
		TableOccupancy:    tableOccupancy,
		PopularDish:       popularProducts,
		NewDish:           newProducts,
		MonthlyOrderCount: monthlyOrderCount,
	}

	// Logging
	repo.log.Info("Fetched dashboard data", zap.Any("summary", summary))

	return summary, nil
}
