package seeder

import "project/domain"

func OrderSeed() []domain.Order {
	return []domain.Order{
		{
			TableID: 1,
			Name:    "John Doe",
			OrderItems: []domain.OrderItem{
				{
					ProductID: 2,
					Quantity:  2,
				},
			},
		},
	}
}
