package domain

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID          uint      `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	CategoryID  int       `gorm:"not null" json:"-" binding:"required,gt=0" form:"category_id" example:"1"`
	Category    Category  `gorm:"foreignKey:CategoryID;references:ID" swaggerignore:"true"`
	Image       string    `gorm:"size:255;not null" json:"image" binding:"omitempty" example:"/image/product.png"`
	Name        string    `gorm:"size:100;unique" json:"name" form:"name"`
	CodeProduct string    `gorm:"size:50;unique" json:"code_product" form:"code_product"`
	Quantity    int       `gorm:"not null" binding:"required,gt=0" json:"quantity" form:"quantity" example:"50"`
	Price       float64   `gorm:"type:decimal(10,2);not null" binding:"required,gt=0" json:"price" form:"price" example:"699.99"`
	Stock       string    `gorm:"size:20;check:stock IN ('In Stock', 'Out Of Stock')" json:"stock" example:"In Stock"`
	Status      string    `gorm:"not null" binding:"required,gt=0" json:"status" form:"status" example:"Active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
}

func (v *Inventory) BeforeSave(tx *gorm.DB) (err error) {
	if v.Quantity > 0 {
		v.Stock = "In Stock"
	} else {
		v.Stock = "Out Of Stock"
	}
	return nil
}

func InventorySeed() []Inventory {
	return []Inventory{
		{
			CategoryID:  1,
			Image:       "/image/coca_cola.png",
			Name:        "Coca Cola",
			CodeProduct: "BEV-001",
			Quantity:    0,
			Price:       1.99,
			Status:      "Active",
		},
		{
			CategoryID:  1,
			Image:       "/image/pepsi.png",
			Name:        "Pepsi",
			CodeProduct: "BEV-002",
			Quantity:    120,
			Price:       1.89,
			Status:      "Active",
		},
		{
			CategoryID:  2,
			Image:       "/image/potato_chips.png",
			Name:        "Potato Chips",
			CodeProduct: "SNK-001",
			Quantity:    80,
			Price:       2.49,
			Status:      "Active",
		},
		{
			CategoryID:  2,
			Image:       "/image/pretzels.png",
			Name:        "Pretzels",
			CodeProduct: "SNK-002",
			Quantity:    60,
			Price:       1.79,
			Status:      "Active",
		},
		{
			CategoryID:  3,
			Image:       "/image/cheesecake.png",
			Name:        "Cheesecake",
			CodeProduct: "DSR-001",
			Quantity:    50,
			Price:       4.99,
			Status:      "Active",
		},
		{
			CategoryID:  3,
			Image:       "/image/ice_cream.png",
			Name:        "Ice Cream",
			CodeProduct: "DSR-002",
			Quantity:    90,
			Price:       3.49,
			Status:      "Active",
		},
		{
			CategoryID:  4,
			Image:       "/image/apple.png",
			Name:        "Apple",
			CodeProduct: "FRT-001",
			Quantity:    150,
			Price:       0.99,
			Status:      "Active",
		},
		{
			CategoryID:  4,
			Image:       "/image/banana.png",
			Name:        "Banana",
			CodeProduct: "FRT-002",
			Quantity:    180,
			Price:       0.79,
			Status:      "Active",
		},
		{
			CategoryID:  5,
			Image:       "/image/carrot.png",
			Name:        "Carrot",
			CodeProduct: "VEG-001",
			Quantity:    200,
			Price:       0.69,
			Status:      "Active",
		},
		{
			CategoryID:  5,
			Image:       "/image/broccoli.png",
			Name:        "Broccoli",
			CodeProduct: "VEG-002",
			Quantity:    100,
			Price:       1.29,
			Status:      "Active",
		},
		{
			CategoryID:  6,
			Image:       "/image/chicken_breast.png",
			Name:        "Chicken Breast",
			CodeProduct: "MEAT-001",
			Quantity:    60,
			Price:       5.99,
			Status:      "Active",
		},
		{
			CategoryID:  6,
			Image:       "/image/beef_steak.png",
			Name:        "Beef Steak",
			CodeProduct: "MEAT-002",
			Quantity:    40,
			Price:       8.99,
			Status:      "Active",
		},
		{
			CategoryID:  7,
			Image:       "/image/milk.png",
			Name:        "Milk",
			CodeProduct: "DAIRY-001",
			Quantity:    200,
			Price:       1.49,
			Status:      "Active",
		},
		{
			CategoryID:  7,
			Image:       "/image/yogurt.png",
			Name:        "Yogurt",
			CodeProduct: "DAIRY-002",
			Quantity:    100,
			Price:       1.99,
			Status:      "Active",
		},
		{
			CategoryID:  8,
			Image:       "/image/bread.png",
			Name:        "Bread",
			CodeProduct: "BAKE-001",
			Quantity:    80,
			Price:       2.29,
			Status:      "Active",
		},
		{
			CategoryID:  8,
			Image:       "/image/croissant.png",
			Name:        "Croissant",
			CodeProduct: "BAKE-002",
			Quantity:    60,
			Price:       1.89,
			Status:      "Active",
		},
		{
			CategoryID:  9,
			Image:       "/image/hot_tea.png",
			Name:        "Hot Tea",
			CodeProduct: "HBEV-001",
			Quantity:    70,
			Price:       1.59,
			Status:      "Active",
		},
		{
			CategoryID:  9,
			Image:       "/image/latte.png",
			Name:        "Latte",
			CodeProduct: "HBEV-002",
			Quantity:    50,
			Price:       2.99,
			Status:      "Active",
		},
		{
			CategoryID:  10,
			Image:       "/image/smoothie.png",
			Name:        "Smoothie",
			CodeProduct: "CBEV-001",
			Quantity:    60,
			Price:       3.49,
			Status:      "Active",
		},
		{
			CategoryID:  10,
			Image:       "/image/orange_juice.png",
			Name:        "Orange Juice",
			CodeProduct: "CBEV-002",
			Quantity:    100,
			Price:       2.49,
			Status:      "Active",
		},
	}
}
