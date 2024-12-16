package domain

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID              uint          `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	TableID         uint          `gorm:"not null" json:"table_id" binding:"required" form:"table_id" example:"1"`
	Table           Table         `gorm:"foreignKey:TableID;references:ID"`
	Name            string        `gorm:"size:100" json:"name" form:"name"`
	CodeOrder       string        `gorm:"size:50;unique" json:"code_order"`
	Tax             float64       `gorm:"type:decimal(4,2);not null;default:10.0" json:"tax"`
	Amount          float64       `gorm:"type:decimal(10,2);not null" binding:"required,gt=0" json:"amount" form:"amount" example:"699.99"`
	PaymentMethodID uint          `gorm:"not null" json:"payment_method_id" binding:"required" form:"payment_method_id" example:"1"`
	PaymentMethod   PaymentMethod `gorm:"foreignKey:PaymentMethodID;references:ID"`
	Status          string        `gorm:"size:20;check:status IN ('In Process', 'Completed', 'Cancelled');default:In Process" json:"status" example:"In Process"`
	CreatedAt       time.Time     `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
	OrderItems      []OrderItem   `gorm:"foreignKey:OrderID;references:ID"`
}

func (o *Order) BeforeSave(tx *gorm.DB) (err error) {
	var orderItems []OrderItem

	if err := tx.Where("order_id = ?", o.ID).Find(&orderItems).Error; err != nil {
		return errors.New("failed to calculate order amount")
	}

	fmt.Printf("Found %d order items for OrderID %d\n", len(orderItems), o.ID)

	totalSubTotal := 0.0
	for _, item := range orderItems {
		fmt.Printf("OrderItem SubTotal: %.2f\n", item.SubTotal)
		totalSubTotal += item.SubTotal
	}

	fmt.Printf("Total SubTotal: %.2f\n", totalSubTotal)

	o.Amount = totalSubTotal + (totalSubTotal * o.Tax / 100)
	fmt.Printf("Order Amount: %.2f\n", o.Amount)

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
			ID:        1,
			TableID:   1,
			Name:      "John Doe",
			CodeOrder: "ORD001",
			// Amount:          150.75,
			PaymentMethodID: 1,
			Status:          "In Process",
			OrderItems: []OrderItem{
				{
					OrderID:   1,
					ProductID: 1,
					Quantity:  2,
					// SubTotal:  300.50,
					Status: "In The Kitchen",
				},
			},
		},
		{
			ID:        2,
			TableID:   2,
			Name:      "Alex",
			CodeOrder: "ORD002",
			// Amount:          245.50,
			PaymentMethodID: 2,
			Status:          "Completed",
		},
		{
			ID:        3,
			TableID:   3,
			Name:      "Elia",
			CodeOrder: "ORD003",
			// Amount:          89.99,
			PaymentMethodID: 3,
			Status:          "Cancelled",
		},
		{
			ID:        4,
			TableID:   4,
			Name:      "Smith",
			CodeOrder: "ORD004",
			// Amount:          175.25,
			PaymentMethodID: 2,
			Status:          "In Process",
		},
		{
			ID:        5,
			TableID:   5,
			Name:      "Bob Brown",
			CodeOrder: "ORD005",
			// Amount:          320.80,
			PaymentMethodID: 1,
			Status:          "Completed",
		},
	}
}

func OrderItemsSeed() []OrderItem {
	return []OrderItem{
		{
			OrderID:   1,
			ProductID: 1,
			Quantity:  2,
			// SubTotal:  300.50,
			Status: "In The Kitchen",
		},
		{
			OrderID:   1,
			ProductID: 2,
			Quantity:  1,
			// SubTotal:  150.25,
			Status: "Cooking Now",
		},
		{
			OrderID:   2,
			ProductID: 3,
			Quantity:  3,
			// SubTotal:  450.75,
			Status: "Ready To Serve",
		},
		{
			OrderID:   3,
			ProductID: 4,
			Quantity:  1,
			// SubTotal:  89.99,
			Status: "In The Kitchen",
		},
		{
			OrderID:   4,
			ProductID: 5,
			Quantity:  5,
			// SubTotal:  875.00,
			Status: "Cooking Now",
		},
	}
}
