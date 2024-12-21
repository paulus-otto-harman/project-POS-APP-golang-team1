package domain

import (

	"time"

	"gorm.io/gorm"
)

type StatusPayment string

const (
	OrderInProcess StatusPayment = "In Process"
	OrderCompleted StatusPayment = "Completed"
	OrderCancelled StatusPayment = "Cancelled"
)

type StatusKitchen string

const (
	OrderInTheKitchen StatusKitchen = "In The Kitchen"
	OrderCookingNow   StatusKitchen = "Cooking Now"
	OrderReadyToServe StatusKitchen = "Ready To Serve"
)

type Order struct {
	ID              uint           `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	TableID         uint           `gorm:"not null" json:"table_id" example:"1"`
	Table           Table          `gorm:"foreignKey:TableID;references:ID"`
	Name            string         `gorm:"size:100" json:"name"`
	CodeOrder       string         `gorm:"size:50;unique" json:"code_order"`
	Tax             float64        `gorm:"type:decimal(4,2);not null;default:10.0" json:"tax"`
	PaymentMethodID *uint          `gorm:"default:null" json:"payment_method_id" example:"1"`
	PaymentMethod   PaymentMethod  `gorm:"foreignKey:PaymentMethodID;references:ID"`
	StatusPayment   StatusPayment  `gorm:"type:status_payment;default:'In Process'" json:"status_payment" example:"In Process"`
	StatusKitchen   StatusKitchen  `gorm:"type:status_kitchen;default:'In The Kitchen'" json:"status_kitchen" example:"In The Kitchen"`
	OrderItems      []OrderItem    `gorm:"foreignKey:OrderID;references:ID"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}


type OrderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	OrderID   uint      `gorm:"not null" json:"order_id" example:"1" swaggerignore:"true"`
	Order     Order     `gorm:"foreignKey:OrderID;references:ID" swaggerignore:"true"`
	ProductID uint      `gorm:"not null" json:"product_id" binding:"required" example:"1"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID" binding:"-" swaggerignore:"true"`
	Quantity  int       `gorm:"not null" json:"quantity" binding:"gt=0" example:"2"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
}
