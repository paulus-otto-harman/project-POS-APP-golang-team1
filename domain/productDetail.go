package domain

type ProductDetail struct {
	ID           int     `json:"id"`
	Image        string  `json:"image"`
	Name         string  `json:"name"`
	CodeProduct  string  `json:"code_product"`
	Stock        int     `json:"stock"`
	Category     string  `json:"category"`
	Price        float64 `json:"price"`
	Availability string  `json:"availability"`
}
