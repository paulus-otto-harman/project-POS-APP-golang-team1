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

func (repo *DashboardRepository) GetDashboard() (*domain.Dashboard, error) {
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

	// Occupied Tables Query
	err = repo.db.Model(&domain.Order{}).
		Where("orders.status_payment = ?", domain.OrderCompleted).
		Joins("JOIN tables ON orders.table_id = tables.id").
		Count(&occupiedTables).Error
	if err != nil {
		repo.log.Error("Error fetching occupied tables", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch occupied tables: %v", err)
	}

	// Calculate occupancy percentage
	occupancyPercentage := (float64(occupiedTables) / float64(totalTables)) * 100

	// Popular Products Query
	var popularProducts []domain.PopularNewResponse
	err = repo.db.Model(&domain.Order{}).
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

	// New Products Query
	var newProducts []domain.PopularNewResponse
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	err = repo.db.Model(&domain.Product{}).
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

	// Set summary
	summary.DailySales = dailySales
	summary.MonthlySales = monthlySales
	summary.TableOccupancy = occupancyPercentage
	summary.PopularDish = popularProducts
	summary.NewDish = newProducts

	// Logging
	repo.log.Info("Fetched popular dishes", zap.Any("dishes", popularProducts))
	repo.log.Info("Fetched new dishes", zap.Any("dishes", newProducts))

	return &summary, nil
}
