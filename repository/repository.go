package repository

import (
	"project/config"
	"project/database"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	Product          ProductRepository
	UserPermission   UserPermissionRepository
	Auth             AuthRepository
	PasswordReset    PasswordResetRepository
	User             UserRepository
	Reservation      ReservationRepository
	Notification     NotificationRepository
	Category         CategoryRepository
	Order            OrderRepository
	UserNotification UserNotificationRepository
	Revenue          RevenueRepository
}

func NewRepository(db *gorm.DB, cacher database.Cacher, config config.Config, log *zap.Logger) Repository {
	return Repository{
		Auth:             *NewAuthRepository(db, cacher, log),
		PasswordReset:    *NewPasswordResetRepository(db, log),
		User:             *NewUserRepository(db, log),
		Reservation:      *NewReservationRepository(db, log),
		Notification:     *NewNotificationRepository(db, log),
		Category:         *NewCategoryRepository(db, log),
		Order:            *NewOrderRepository(db, log),
		UserNotification: *NewUserNotificationRepository(db, log),
		Product:          *NewProductRepository(db, log),
		UserPermission:   *NewUserPermissionRepository(db, log),
		Revenue:          *NewRevenueRepository(db, log),
	}
}
