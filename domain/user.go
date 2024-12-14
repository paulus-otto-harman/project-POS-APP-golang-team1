package domain

import (
	"gorm.io/gorm"
	"time"
)

type UserRole string

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
