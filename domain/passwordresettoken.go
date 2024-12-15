package domain

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type PasswordResetToken struct {
	ID        uuid.UUID    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email     string       `gorm:"size:30" json:"email" json:"email"`
	Otp       string       `gorm:"size:8" json:"otp" json:"otp"`
	CreatedAt time.Time    `gorm:"default:now()" json:"created_at"`
	ExpiredAt sql.NullTime `json:"expired_at"`
}
