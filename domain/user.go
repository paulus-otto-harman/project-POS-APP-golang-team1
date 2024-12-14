package domain

import (
	"gorm.io/gorm"
	"time"
)

type UserRole string

const (
	SuperAdmin UserRole = "super admin"
	Admin      UserRole = "admin"
	Staff      UserRole = "staff"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"unique" example:"admin@mail.com" json:"email"`
	Password  string         `example:"password" json:"password"`
	Role      UserRole       `gorm:"type:user_role" json:"role"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Profile             Profile              `json:"profile"`
	Permissions         []Permission         `gorm:"many2many:user_permissions;" json:"permissions"`
	PasswordResetTokens []PasswordResetToken `gorm:"foreignKey:Email;references:Email" json:"-"`
}

func UserSeed() []User {
	return []User{
		{
			Email:    "super@mail.com",
			Password: "super",
			Role:     SuperAdmin,
			Profile: Profile{
				FullName: "Super Admin",
				Phone:    "00",
				Salary:   150,
			},
		},
		{
			Email:    "admin@mail.com",
			Password: "admin",
			Role:     Admin,
			Profile: Profile{
				FullName: "Admin Satu",
				Phone:    "01",
				Salary:   100,
			},
			Permissions: []Permission{
				{ID: 1, Name: "Dashboard"},
				{ID: 2, Name: "Reports"},
				{ID: 6, Name: "Settings"},
			},
		},
		{
			Email:    "staff@mail.com",
			Password: "staff",
			Role:     Staff,
			Profile: Profile{
				FullName: "Staff Satu",
				Phone:    "02",
				Salary:   50,
			},
			Permissions: []Permission{
				{ID: 4, Name: "Orders"},
				{ID: 5, Name: "Customers"},
			},
		},
	}
}
