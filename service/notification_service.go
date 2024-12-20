package service

import (
	"project/domain"
	"project/repository"
	"time"

	"go.uber.org/zap"
)

type NotificationService interface {
	CreateNotificationLowStock() error
	All(userID uint, status string) ([]domain.Notification, error)
	Update(id uint, status string) error
	Delete(id uint) error
	BatchUpdate(ids []uint, status string) error
}

type notificationService struct {
	repo repository.Repository
	log  *zap.Logger
}

// CreateNotification implements NotificationService.
func (n *notificationService) CreateNotificationLowStock() error {
	products, _, err := n.repo.Product.All(0, 0, "", "", "Low Stock", 0, 0.0, 0.0)
	if err != nil {
		n.log.Error("Failed to fetch products", zap.Error(err))
		return err
	}

	n.log.Info("Get admin data from users")
	admins, err := n.repo.User.GetByRole("admin")
	if err != nil {
		n.log.Error("Failed to fetch admins", zap.Error(err))
		return err
	}

	for _, product := range products {
		// Create a new notification
		newNotif := domain.Notification{
			Title:        "Low Inventory Alert",
			Content:      "This is to notify you that the following items are running low in stock:",
			ProductName:  product.Name,
			ProductImage: product.Image,
		}

		err = n.repo.Notification.Create(&newNotif)
		if err != nil {
			return err
		}

		// The ID is now set automatically after Create
		createdNotifID := newNotif.ID

		// Now, create UserNotification for each admin
		for _, admin := range admins {
			// Create a UserNotification for each admin
			userNotif := domain.UserNotification{
				UserID:         admin.ID,
				NotificationID: createdNotifID, // Use the ID from the created notification
				Status:         "unread",       // Default status
				CreatedAt:      time.Now(),     // Or use the default time set in BeforeCreate
			}

			// Insert the UserNotification into the database
			err := n.repo.UserNotification.Create(userNotif)
			if err != nil {
				return err
			}
			n.log.Info("Created UserNotification for admin", zap.String("UserID", admin.FullName), zap.Uint("NotificationID", createdNotifID))
		}
	}
	return nil
}

// Delete implements NotificationService.
func (n *notificationService) Delete(id uint) error {
	n.log.Info("Deleting a notification")
	return n.repo.Notification.Delete(id)
}

// All implements NotificationService.
func (n *notificationService) All(userID uint, status string) ([]domain.Notification, error) {
	n.log.Info("Fetching all notifications")
	return n.repo.Notification.All(userID, status)
}

// UpdateNotification implements NotificationService.
func (n *notificationService) Update(id uint, status string) error {
	n.log.Info("Updating a notification")
	return n.repo.Notification.Update(id, status)
}

func (n *notificationService) BatchUpdate(ids []uint, status string) error {
	n.log.Info("Batch updating notifications")
	return n.repo.Notification.BatchUpdate(ids, status)
}

func NewNotificationService(repo repository.Repository, log *zap.Logger) NotificationService {
	return &notificationService{
		repo: repo,
		log:  log,
	}
}
