package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           uint      `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	CategoryID   int       `gorm:"not null" json:"-" binding:"required,gt=0" form:"category_id" example:"1"`
	Category     Category  `gorm:"foreignKey:CategoryID;references:ID" swaggerignore:"true"`
	Image        string    `gorm:"size:255;not null" json:"image" binding:"omitempty" example:"/image/product.png"`
	Name         string    `gorm:"size:100;unique" json:"name" form:"name"`
	CodeProduct  string    `gorm:"size:50;unique" json:"code_product" form:"code_product"`
	Stock        int       `gorm:"not null" binding:"required,gt=0" json:"stock" form:"stock" example:"50"`
	Price        float64   `gorm:"type:decimal(10,2);not null" binding:"required,gt=0" json:"price" form:"price" example:"699.99"`
	Availability string    `gorm:"size:20;check:availability IN ('In Stock', 'Out Of Stock')" json:"availability" example:"In Stock"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
}

func (v *Product) BeforeSave(tx *gorm.DB) (err error) {
	if v.Stock > 0 {
		v.Availability = "In Stock"
	} else {
		v.Availability = "Out Of Stock"
	}
	return nil
}

func ProductSeed() []Product {
	return []Product{
		{
			CategoryID:  1,
			Image:       "/image/coca_cola.png",
			Name:        "Coca Cola",
			CodeProduct: "BEV-001",
			Stock:       0,
			Price:       1.99,
		},
		{
			CategoryID:  1,
			Image:       "/image/pepsi.png",
			Name:        "Pepsi",
			CodeProduct: "BEV-002",
			Stock:       120,
			Price:       1.89,
		},
		{
			CategoryID:  2,
			Image:       "/image/potato_chips.png",
			Name:        "Potato Chips",
			CodeProduct: "SNK-001",
			Stock:       80,
			Price:       2.49,
		},
		{
			CategoryID:  2,
			Image:       "/image/pretzels.png",
			Name:        "Pretzels",
			CodeProduct: "SNK-002",
			Stock:       60,
			Price:       1.79,
		},
		{
			CategoryID:  3,
			Image:       "/image/cheesecake.png",
			Name:        "Cheesecake",
			CodeProduct: "DSR-001",
			Stock:       50,
			Price:       4.99,
		},
		{
			CategoryID:  3,
			Image:       "/image/ice_cream.png",
			Name:        "Ice Cream",
			CodeProduct: "DSR-002",
			Stock:       90,
			Price:       3.49,
		},
		{
			CategoryID:  4,
			Image:       "/image/apple.png",
			Name:        "Apple",
			CodeProduct: "FRT-001",
			Stock:       150,
			Price:       0.99,
		},
		{
			CategoryID:  4,
			Image:       "/image/banana.png",
			Name:        "Banana",
			CodeProduct: "FRT-002",
			Stock:       180,
			Price:       0.79,
		},
		{
			CategoryID:  5,
			Image:       "/image/carrot.png",
			Name:        "Carrot",
			CodeProduct: "VEG-001",
			Stock:       200,
			Price:       0.69,
		},
		{
			CategoryID:  5,
			Image:       "/image/broccoli.png",
			Name:        "Broccoli",
			CodeProduct: "VEG-002",
			Stock:       100,
			Price:       1.29,
		},
		{
			CategoryID:  6,
			Image:       "/image/chicken_breast.png",
			Name:        "Chicken Breast",
			CodeProduct: "MEAT-001",
			Stock:       60,
			Price:       5.99,
		},
		{
			CategoryID:  6,
			Image:       "/image/beef_steak.png",
			Name:        "Beef Steak",
			CodeProduct: "MEAT-002",
			Stock:       40,
			Price:       8.99,
		},
		{
			CategoryID:  7,
			Image:       "/image/milk.png",
			Name:        "Milk",
			CodeProduct: "DAIRY-001",
			Stock:       200,
			Price:       1.49,
		},
		{
			CategoryID:  7,
			Image:       "/image/yogurt.png",
			Name:        "Yogurt",
			CodeProduct: "DAIRY-002",
			Stock:       100,
			Price:       1.99,
		},
		{
			CategoryID:  8,
			Image:       "/image/bread.png",
			Name:        "Bread",
			CodeProduct: "BAKE-001",
			Stock:       80,
			Price:       2.29,
		},
		{
			CategoryID:  8,
			Image:       "/image/croissant.png",
			Name:        "Croissant",
			CodeProduct: "BAKE-002",
			Stock:       60,
			Price:       1.89,
		},
		{
			CategoryID:  9,
			Image:       "/image/hot_tea.png",
			Name:        "Hot Tea",
			CodeProduct: "HBEV-001",
			Stock:       70,
			Price:       1.59,
		},
		{
			CategoryID:  9,
			Image:       "/image/latte.png",
			Name:        "Latte",
			CodeProduct: "HBEV-002",
			Stock:       50,
			Price:       2.99,
		},
		{
			CategoryID:  10,
			Image:       "/image/smoothie.png",
			Name:        "Smoothie",
			CodeProduct: "CBEV-001",
			Stock:       60,
			Price:       3.49,
		},
		{
			CategoryID:  10,
			Image:       "/image/orange_juice.png",
			Name:        "Orange Juice",
			CodeProduct: "CBEV-002",
			Stock:       100,
			Price:       2.49,
		},
	}
}
