package domain

import (
	"time"
)

type BestSeller struct {
	ID           uint      `gorm:"primaryKey" json:"id" `
	ProductID    uint      `form:"product_id" json:"product_id"`
	Product      Product   `gorm:"foreignKey:ProductID" json:"product"`
	Date         time.Time `gorm:"type:date" json:"date"`
	SellPrice    float64   `json:"sell_price"`
	Profit       float64   `json:"profit"`
	ProfitMargin float32   `json:"profit_margin"`
	Revenue      float64   `json:"revenue"`
}
