package seeder

import "project/domain"

func CategorySeed() []domain.Category {
	return []domain.Category{
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
