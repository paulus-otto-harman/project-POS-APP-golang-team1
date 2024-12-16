package domain

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID              uint           `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	TableID         uint           `gorm:"not null" json:"table_id" binding:"required" form:"table_id" example:"1"`
	Table           Table          `gorm:"foreignKey:TableID;references:ID"`
	Name            string         `gorm:"size:100" json:"name" form:"name"`
	CodeOrder       string         `gorm:"size:50;unique" json:"code_order"`
	Tax             float64        `gorm:"type:decimal(4,2);not null;default:10.0" json:"tax"`
	Amount          float64        `gorm:"type:decimal(10,2);not null" binding:"required,gt=0" json:"amount" form:"amount" example:"699.99"`
	PaymentMethodID uint           `gorm:"not null" json:"payment_method_id" binding:"required" form:"payment_method_id" example:"1"`
	PaymentMethod   PaymentMethod  `gorm:"foreignKey:PaymentMethodID;references:ID"`
	Status          string         `gorm:"size:20;check:status IN ('In Process', 'Completed', 'Cancelled');default:In Process" json:"status" example:"In Process"`
	OrderItems      []OrderItem    `gorm:"foreignKey:OrderID;references:ID"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}

func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
	var totalSubTotal float64

	for i, item := range o.OrderItems {
		var product Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			return fmt.Errorf("product not found for product_id %d", item.ProductID)
		}
		o.OrderItems[i].SubTotal = product.Price * float64(item.Quantity)
		totalSubTotal += o.OrderItems[i].SubTotal
	}

	o.Amount = totalSubTotal + (totalSubTotal * o.Tax / 100)
	return nil
}

type OrderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	OrderID   uint      `gorm:"not null" json:"order_id" binding:"required,gt=0" form:"order_id" example:"1"`
	Order     Order     `gorm:"foreignKey:OrderID;references:ID"`
	ProductID uint      `gorm:"not null" json:"product_id" binding:"required,gt=0" form:"product_id" example:"1"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID"`
	Quantity  int       `gorm:"not null" json:"quantity" binding:"required,gt=0" form:"quantity" example:"2"`
	SubTotal  float64   `gorm:"type:decimal(10,2);not null" binding:"required,gt=0" json:"sub_total" form:"sub_total" example:"699.99"`
	Status    string    `gorm:"size:20;check:status IN ('In The Kitchen', 'Cooking Now', 'Ready To Serve');default:In The Kitchen" json:"status" example:"In The Kitchen"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
}

func (oi *OrderItem) BeforeSave(tx *gorm.DB) (err error) {

	if oi.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	var product Product
	if err := tx.First(&product, oi.ProductID).Error; err != nil {
		return errors.New("product not found")
	}

	oi.SubTotal = product.Price * float64(oi.Quantity)
	return nil
}

func OrderSeed() []Order {
	return []Order{
		{
			ID:              1,
			TableID:         1,
			Name:            "John Doe",
			CodeOrder:       "ORD001",
			Tax:             10.0,
			PaymentMethodID: 1,
			Status:          "In Process",
			OrderItems: []OrderItem{
				{
					ProductID: 1,
					Quantity:  2,
					Status:    "In The Kitchen",
				},
				{
					ProductID: 2,
					Quantity:  1,
					Status:    "Cooking Now",
				},
			},
		},
		{
			ID:              2,
			TableID:         2,
			Name:            "Alex",
			CodeOrder:       "ORD002",
			Tax:             10.0,
			PaymentMethodID: 2,
			Status:          "Completed",
			OrderItems: []OrderItem{
				{
					ProductID: 3,
					Quantity:  3,
					Status:    "Ready To Serve",
				},
				{
					ProductID: 12,
					Quantity:  3,
					Status:    "Ready To Serve",
				},
				{
					ProductID: 6,
					Quantity:  3,
					Status:    "Ready To Serve",
				},
			},
		},
		{
			ID:              3,
			TableID:         3,
			Name:            "Elia",
			CodeOrder:       "ORD003",
			Tax:             10.0,
			PaymentMethodID: 3,
			Status:          "Cancelled",
			OrderItems: []OrderItem{
				{
					ProductID: 4,
					Quantity:  1,
					Status:    "In The Kitchen",
				},
				{
					ProductID: 15,
					Quantity:  1,
					Status:    "In The Kitchen",
				},
			},
		},
		{
			ID:              4,
			TableID:         4,
			Name:            "Smith",
			CodeOrder:       "ORD004",
			Tax:             10.0,
			PaymentMethodID: 2,
			Status:          "In Process",
			OrderItems: []OrderItem{
				{
					ProductID: 5,
					Quantity:  5,
					Status:    "Cooking Now",
				},
				{
					ProductID: 17,
					Quantity:  5,
					Status:    "Cooking Now",
				},
				{
					ProductID: 2,
					Quantity:  5,
					Status:    "Cooking Now",
				},
			},
		},
		{
			ID:              5,
			TableID:         5,
			Name:            "Bob Brown",
			CodeOrder:       "ORD005",
			Tax:             10.0,
			PaymentMethodID: 1,
			Status:          "Completed",
			OrderItems: []OrderItem{
				{
					ProductID: 5,
					Quantity:  5,
					Status:    "Cooking Now",
				},
			},
		},
	}

}
