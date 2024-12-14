package service

import (
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type NotificationService interface {
	CreateNotificationLowStock() error
	GetAllNotifications(status string) ([]domain.Notification, error)
	UpdateNotification(id uint, status string) error
	DeleteNotification(id uint) error
	BatchUpdateNotifications(ids []uint, status string) error
}

type notificationService struct {
	repo repository.NotificationRepository
	log  *zap.Logger
}

// GetAllNotifications implements NotificationService.
func (n *notificationService) GetAllNotifications(status string) ([]domain.Notification, error) {
	n.log.Info("Fetching all notifications")
	return n.repo.GetAll(status)
}

// UpdateNotification implements NotificationService.
func (n *notificationService) UpdateNotification(id uint, status string) error {
	n.log.Info("Updating a notification")
	return n.repo.Update(id, status)
}

func (n *notificationService) BatchUpdateNotifications(ids []uint, status string) error {
	n.log.Info("Batch updating notifications")
	return n.repo.BatchUpdate(ids, status)
}

func NewNotificationService(repo repository.NotificationRepository, log *zap.Logger) NotificationService {
	return &notificationService{
		repo: repo,
		log:  log,
	}
}
