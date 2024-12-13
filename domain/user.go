package domain

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	FullName  string         `json:"full_name"`
	Email     string         `gorm:"unique" example:"admin@mail.com" json:"email"`
	Password  string         `example:"password" json:"password"`
	Role      string         `gorm:"default:staff" json:"role"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	PasswordResetTokens []PasswordResetToken `gorm:"foreignKey:Email;references:Email" json:"-"`
}

func UserSeed() []User {
	return []User{
		{
			FullName: "Super Admin",
			Email:    "admin@mail.com",
			Password: "admin",
			Role:     "admin",
		},
		{
			FullName: "Staf Satu",
			Email:    "staf@mail.com",
			Password: "staf",
			Role:     "staff",
		},
	}
}
