package seeder

import (
	"project/domain"
	"project/helper"
)

func User() []domain.User {
	return []domain.User{
		{
			Email:       "super@mail.com",
			Password:    helper.HashPassword("super"),
			Role:        domain.SuperAdmin,
			FullName:    "Super Admin",
			PhoneNumber: "00",
			Salary:      150,
		},
		{
			Email:       "admin@mail.com",
			Password:    helper.HashPassword("admin"),
			Role:        domain.Admin,
			FullName:    "Admin Satu",
			PhoneNumber: "01",
			Salary:      100,
			Permissions: []domain.Permission{
				{ID: 1, Name: "Dashboard"},
				{ID: 2, Name: "Reports"},
				{ID: 6, Name: "Settings"},
			},
		},
		{
			Email:       "staff@mail.com",
			Password:    helper.HashPassword("staff"),
			Role:        domain.Staff,
			FullName:    "Staff Satu",
			PhoneNumber: "02",
			Salary:      50,
			Permissions: []domain.Permission{
				{ID: 4, Name: "Orders"},
				{ID: 5, Name: "Customers"},
			},
		},
	}
}
