package domain

type Dashboard struct {
	DailySales     float64              `json:"daily_sales"`
	MonthlySales   float64              `json:"monthly_sales"`
	TableOccupancy float64              `json:"table_occupancy"`
	PopularDish    []PopularNewResponse `json:"popular_dish"`
	NewDish        []PopularNewResponse `json:"new_dish"`
}

type PopularNewResponse struct {
	Name         string  `json:"name"`
	OrderCount   int     `json:"order_count"`
	Image        string  `json:"image"`
	Availability string  `json:"availability"`
	Price        float64 `json:"price"`
}
