package service

import (
	"project/repository"

	"go.uber.org/zap"
)

type Service struct {
	Auth          AuthService
	PasswordReset PasswordResetService
	User          UserService
	Reservation   ReservationService
	Notification  NotificationService
	Category      CategoryService
	Product       ProductService
}

func NewService(repo repository.Repository, log *zap.Logger) Service {
	return Service{
		Auth:          NewAuthService(repo.Auth, log),
		PasswordReset: NewPasswordResetService(repo.PasswordReset, log),
		User:          NewUserService(repo.User, log),
		Notification:  NewNotificationService(repo, log),
		Reservation:   NewReservationService(repo.Reservation, log),
		Category:      NewCategoryService(repo.Category, log),
		Product:       NewProductService(repo.Product, log),
	}
}
