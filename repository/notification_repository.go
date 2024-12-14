package repository

import (
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewNotificationRepository(db *gorm.DB, log *zap.Logger) *NotificationRepository {
	return &NotificationRepository{db: db, log: log}
}

func (repo NotificationRepository) Create(notification domain.Notification) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&notification).Error; err != nil {
			repo.log.Error("Error creating notification", zap.Error(err))
			return err
		}
		return nil
	})
}

func (repo NotificationRepository) GetAll(status string) ([]domain.Notification, error) {
	var notifications []domain.Notification
	query := repo.db
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&notifications).Error; err != nil {
		repo.log.Error("Error fetching notifications", zap.Error(err))
		return nil, err
	}
	return notifications, nil
}

func (repo NotificationRepository) Update(id uint, status string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Notification{}).Where("id = ?", id).Update("status", status).Error; err != nil {
			repo.log.Error("Error updating notification status", zap.Error(err))
			return err
		}
		return nil
	})
}

func (repo NotificationRepository) Delete(id uint) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&domain.Notification{}, id).Error; err != nil {
			repo.log.Error("Error deleting notification", zap.Error(err))
			return err
		}
		return nil
	})
}

func (repo NotificationRepository) BatchUpdate(ids []uint, status string) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		// Perform the batch update
		if err := tx.Model(&domain.Notification{}).
			Where("id IN ?", ids).
			Update("status", status).Error; err != nil {
			repo.log.Error("Error batch updating notification statuses", zap.Error(err))
			return err
		}
		return nil
	})
}
