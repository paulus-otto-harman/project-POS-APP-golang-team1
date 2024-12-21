package domain

import (
	"time"

)

type Product struct {
	ID           uint       `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	CategoryID   int        `gorm:"not null" json:"category_id" binding:"required,gt=0" form:"category_id" example:"1"`
	Category     Category   `gorm:"foreignKey:CategoryID;references:ID" swaggerignore:"true"`
	Image        string     `gorm:"size:255;not null" json:"image" binding:"required" example:"/image/product.png"`
	Name         string     `gorm:"size:100;unique" json:"name" form:"name"`
	CodeProduct  string     `gorm:"size:50;unique" json:"code_product" form:"code_product"`
	Stock        int        `gorm:"not null" binding:"required,gt=0" json:"stock" form:"stock" example:"50"`
	Price        float64    `gorm:"type:decimal(10,2);not null" binding:"required,gt=0" json:"price" form:"price" example:"699.99"`
	Availability string     `gorm:"size:20;check:availability IN ('In Stock', 'Low Stock', 'Out Of Stock')" json:"availability" example:"In Stock"`
	Status       string     `gorm:"not null;default:Active;check:status IN ('Active', 'Inactive')" binding:"required,gt=0" json:"status" form:"status" example:"Active"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

type ProductRevenue struct {
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	TotalRevenue float64 `json:"total_revenue"`
}