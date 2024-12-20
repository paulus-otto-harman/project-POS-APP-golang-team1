package service

import (
	"project/config"
	"project/repository"

	"go.uber.org/zap"
)

type Service struct {
	Auth          AuthService
	Email         EmailService
	Otp           OtpService
	PasswordReset PasswordResetService
	User          UserService
	Reservation   ReservationService
	Notification  NotificationService
	Category      CategoryService
	Product       ProductService
	Order         OrderService
	Dashboard     DashboardService
}

func NewService(repo repository.Repository, appConfig config.Config, log *zap.Logger) Service {
	return Service{
		Auth:          NewAuthService(repo.User, log),
		Email:         NewEmailService(appConfig.Email, log),
		Otp:           NewOtpService(log),
		PasswordReset: NewPasswordResetService(repo.PasswordReset, log),
		User:          NewUserService(repo, log),
		Notification:  NewNotificationService(repo, log),
		Reservation:   NewReservationService(repo.Reservation, log),
		Category:      NewCategoryService(repo.Category, log),
		Product:       NewProductService(repo.Product, log),
		Order:         NewOrderService(repo.Order, log),
		Dashboard:     NewDashboardService(repo.Dashboard, log),
	}
}
