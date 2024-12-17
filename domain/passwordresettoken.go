package domain

import (
	"github.com/google/uuid"
	"time"
)

type PasswordResetToken struct {
	Token     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"token"`
	Email     string    `json:"email" json:"email"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	ExpiredAt time.Time `gorm:"default:now() + '5 minutes'::interval" json:"expired_at"`
}
