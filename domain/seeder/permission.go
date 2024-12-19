package seeder

import "project/domain"

func Permission() []domain.Permission {
	return []domain.Permission{
		{Name: "Dashboard"},
		{Name: "Reports"},
		{Name: "Inventory"},
		{Name: "Orders"},
		{Name: "Customers"},
		{Name: "Settings"},
	}
}
