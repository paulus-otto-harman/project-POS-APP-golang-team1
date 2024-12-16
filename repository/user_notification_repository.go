package repository

import (
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserNotificationRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewUserNotificationRepository(db *gorm.DB, log *zap.Logger) *UserNotificationRepository {
	return &UserNotificationRepository{db: db, log: log}
}

func (repo UserNotificationRepository) Create(userNotifInput domain.UserNotification) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&userNotifInput).Error; err != nil {
			repo.log.Error("Error creating user notification", zap.Error(err))
			return err
		}
		return nil
	})
}
