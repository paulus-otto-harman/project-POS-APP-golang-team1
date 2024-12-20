package domain

type Dashboard struct {
	DailySales     float64 `json:"daily_sales"`
	MonthlySales   float64 `json:"monthly_sales"`
	TableOccupancy float64 `json:"table_occupancy"`
}
