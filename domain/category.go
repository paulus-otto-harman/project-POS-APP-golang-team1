package domain

import (
	"time"
)

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	Icon        string    `gorm:"size:255;not null" json:"icon,omitempty" binding:"omitempty" example:"/icon/category.png"`
	Name        string    `gorm:"size:100;unique" json:"name" form:"name"`
	Description string    `gorm:"type:text" example:"lorem" json:"description,omitempty" form:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
}

func CategorySeed() []Category {
	return []Category{
		{
			Icon:        "/icon/beverage.png",
			Name:        "Beverage",
			Description: "All kinds of beverages including soft drinks, coffee, and tea",
		},
		{
			Icon:        "/icon/snack.png",
			Name:        "Snacks",
			Description: "Light snacks and finger foods",
		},
		{
			Icon:        "/icon/dessert.png",
			Name:        "Desserts",
			Description: "Sweet dishes like cakes, pastries, and ice creams",
		},
		{
			Icon:        "/icon/fruit.png",
			Name:        "Fruits",
			Description: "Fresh and seasonal fruits",
		},
		{
			Icon:        "/icon/vegetable.png",
			Name:        "Vegetables",
			Description: "Fresh vegetables for healthy meals",
		},
		{
			Icon:        "/icon/meat.png",
			Name:        "Meat",
			Description: "All types of fresh and processed meat",
		},
		{
			Icon:        "/icon/dairy.png",
			Name:        "Dairy",
			Description: "Milk, cheese, yogurt, and other dairy products",
		},
		{
			Icon:        "/icon/bakery.png",
			Name:        "Bakery",
			Description: "Freshly baked bread, buns, and cakes",
		},
		{
			Icon:        "/icon/beverages_hot.png",
			Name:        "Hot Beverages",
			Description: "Coffee, tea, and other hot drinks",
		},
		{
			Icon:        "/icon/beverages_cold.png",
			Name:        "Cold Beverages",
			Description: "Chilled drinks including soda, juice, and smoothies",
		},
	}
}
