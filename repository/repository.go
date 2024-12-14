package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project/config"
	"project/database"
)

type Repository struct {
	Auth          AuthRepository
	PasswordReset PasswordResetRepository
	User          UserRepository
	Notification  NotificationRepository
	Category      CategoryRepository
}

func NewRepository(db *gorm.DB, cacher database.Cacher, config config.Config, log *zap.Logger) Repository {
	return Repository{
		Auth:          *NewAuthRepository(db, cacher, log),
		PasswordReset: *NewPasswordResetRepository(db, log),
		User:          *NewUserRepository(db, log),
		Notification:  *NewNotificationRepository(db, log),
		Category: *NewCategoryRepository(db, log),
	}
}
