package seeder

import (
	"project/domain"
	"project/helper"
)

const (
	SuperAdmin domain.UserRole = "super admin"
	Admin      domain.UserRole = "admin"
	Staff      domain.UserRole = "staff"
)

func User() []domain.User {
	return []domain.User{
		{
			Email:    "super@mail.com",
			Password: helper.HashPassword("super"),
			Role:     SuperAdmin,
			Profile: domain.Profile{
				FullName: "Super Admin",
				Phone:    "00",
				Salary:   150,
			},
		},
		{
			Email:    "admin@mail.com",
			Password: helper.HashPassword("admin"),
			Role:     Admin,
			Profile: domain.Profile{
				FullName: "Admin Satu",
				Phone:    "01",
				Salary:   100,
			},
			Permissions: []domain.Permission{
				{ID: 1, Name: "Dashboard"},
				{ID: 2, Name: "Reports"},
				{ID: 6, Name: "Settings"},
			},
		},
		{
			Email:    "staff@mail.com",
			Password: helper.HashPassword("staff"),
			Role:     Staff,
			Profile: domain.Profile{
				FullName: "Staff Satu",
				Phone:    "02",
				Salary:   50,
			},
			Permissions: []domain.Permission{
				{ID: 4, Name: "Orders"},
				{ID: 5, Name: "Customers"},
			},
		},
	}
}
