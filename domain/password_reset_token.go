package domain

import (
	"github.com/google/uuid"
	"time"
)

type PasswordResetToken struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID          uint
	User            User
	Email           string     `gorm:"size:30" json:"email" json:"email"`
	Otp             string     `gorm:"size:8" json:"otp" json:"otp"`
	CreatedAt       time.Time  `gorm:"default:now()" json:"created_at"`
	ExpiredAt       time.Time  `gorm:"default:now() + '3 minutes'::interval" json:"expired_at"`
	ValidatedAt     *time.Time `json:"validated_at"`
	PasswordResetAt *time.Time `json:"password_reset_at"`
}
