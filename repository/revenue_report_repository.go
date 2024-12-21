package repository

import (
	"fmt"
	"project/domain"
	"strings"

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

func (repo *RevenueRepository) GetTotalRevenueByStatus() (map[string]interface{}, error) {
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

func (repo *RevenueRepository) GetMonthlyRevenue(statusPayment string, year int) (map[string]float64, error) {
	result := make(map[string]float64)

	// Initialize all months with zero values
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	for _, month := range months {
		result[month] = 0
	}

	type MonthlyRevenue struct {
		Month   string  `json:"month"`
		Revenue float64 `json:"revenue"`
	}

	var revenues []MonthlyRevenue

	query := repo.db.Table("order_details").
		Select("TO_CHAR(TO_DATE(date_order, 'FMDay, dd-Mon-yyyy'), 'FMMonth') AS month, SUM(total) AS revenue").
		Where("EXTRACT(YEAR FROM TO_DATE(date_order, 'FMDay, dd-Mon-yyyy')) = ?", year).
		Group("TO_CHAR(TO_DATE(date_order, 'FMDay, dd-Mon-yyyy'), 'FMMonth'), EXTRACT(MONTH FROM TO_DATE(date_order, 'FMDay, dd-Mon-yyyy'))").
		Order("EXTRACT(MONTH FROM TO_DATE(date_order, 'FMDay, dd-Mon-yyyy'))")

	if statusPayment != "" {
		query = query.Where("status_payment = ?", statusPayment)
	}

	if err := query.Scan(&revenues).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch revenue data: %w", err)
	}

	// Update result map with actual values
	for _, rev := range revenues {
		month := strings.TrimSpace(rev.Month)
		result[month] = rev.Revenue
	}

	return result, nil
}

func (repo *RevenueRepository) GetProductRevenueDetails() ([]*domain.ProductRevenue, error) {
	var products []*domain.ProductRevenue

	return products, nil
}

func (repo *RevenueRepository) AddDailyBestSeller(profitMargin float64) {
	repo.db.Exec(`
		INSERT INTO best_sellers (product_id, sell_price, profit, profit_margin, revenue)
		SELECT best.product_id, products.price, ` + fmt.Sprintf("%f", profitMargin) + `/100 * best.quantities * products.price, ` + fmt.Sprintf("%f", profitMargin) + `, best.quantities * products.price
			FROM
				(SELECT product_id, SUM(quantity) AS quantities
				FROM order_items
				JOIN orders ON order_items.order_id=orders.id
				WHERE orders.status_payment='Completed'
				AND orders.created_at::DATE=TIMESTAMP 'yesterday'::DATE
				GROUP BY product_id
				ORDER BY 2 DESC
				LIMIT 1) best
			JOIN products ON best.product_id=products.id
	`)
}
