package domain

type Dashboard struct {
	DailySales   float64 `json:"daily_sales"`
	MonthlySales float64 `json:"monthly_sales"`
	TotalTables  int     `json:"total_tables"`
}
