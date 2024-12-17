package domain

import (
	"time"

	"github.com/google/uuid"
)

type PasswordResetToken struct {
	Token     uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"token"`
	Email     string    `gorm:"email" json:"email"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	ExpiredAt time.Time `gorm:"default:now() + '5 minutes'::interval" json:"expired_at"`
}
