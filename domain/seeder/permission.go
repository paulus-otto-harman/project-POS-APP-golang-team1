package seeder

import "project/domain"

func Permission() []domain.Permission {
	return []domain.Permission{
		{Name: "Dashboard"},
		{Name: "Menu"},
		{Name: "Staff"},
		{Name: "Inventory"},
		{Name: "Reports"},
		{Name: "Orders"},
		{Name: "Reservations"},
	}
}
