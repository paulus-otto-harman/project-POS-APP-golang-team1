package seeder

import (
	"project/domain"
	"project/helper"
	"time"
)

func PasswordResetTokenSeed() []domain.PasswordResetToken {
	return []domain.PasswordResetToken{
		{
			UserID:    2,
			Email:     "staff@mail.com",
			Otp:       "111111",
			CreatedAt: time.Now(),
			ExpiredAt: time.Now().Add(time.Hour),
		},
		{
			UserID:    2,
			Email:     "staff@mail.com",
			Otp:       "222222",
			CreatedAt: time.Now(),
			ExpiredAt: time.Now().Add(time.Hour),
		},
		{
			UserID:      2,
			Email:       "staff@mail.com",
			Otp:         "333333",
			CreatedAt:   time.Now(),
			ExpiredAt:   time.Now().Add(time.Hour),
			ValidatedAt: helper.Ptr(time.Now()),
		},
		{
			UserID:          2,
			Email:           "staff@mail.com",
			Otp:             "444444",
			CreatedAt:       time.Now(),
			ExpiredAt:       time.Now().Add(time.Hour),
			ValidatedAt:     helper.Ptr(time.Now()),
			PasswordResetAt: helper.Ptr(time.Now()),
		},
	}
}
