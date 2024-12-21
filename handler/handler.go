package handler

import (
	"project/database"
	"project/domain"
	"project/infra/jwt"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	AuthHandler           AuthController
	PasswordResetHandler  PasswordResetController
	ProfileHandler        ProfileController
	UserHandler           UserController
	ReservationHandler    ReservationController
	NotificationHandler   NotificationController
	CategoryHandler       CategoryController
	ProductHandler        ProductController
	OrderHandler          OrderController
  DashboardHandler     DashboardController
	UserPermissionHandler UserPermissionController
	RevenueHandler        RevenueController
}

func NewHandler(service service.Service, logger *zap.Logger, rdb database.Cacher, jwt jwt.JWT) *Handler {
	return &Handler{
		AuthHandler:           *NewAuthController(service.Auth, logger, rdb, jwt),
		PasswordResetHandler:  *NewPasswordResetController(service, logger),
		UserHandler:           *NewUserController(service, logger),
		ProfileHandler:        *NewProfileController(service, logger, rdb, jwt),
		ReservationHandler:    *NewReservationController(service.Reservation, logger),
		NotificationHandler:   *NewNotificationController(service, logger),
		CategoryHandler:       *NewCategoryController(service.Category, logger),
		ProductHandler:        *NewProductController(service.Product, logger),
		OrderHandler:          *NewOrderController(service.Order, logger),
    DashboardHandler:     *NewDashboardController(service.Dashboard, logger),
		UserPermissionHandler: *NewUserPermissionController(service.UserPermission, logger),
		RevenueHandler:        *NewRevenueController(service.Revenue, logger),
	}
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func BadResponse(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, Response{
		Status:  false,
		Message: message,
	})
}

func GoodResponseWithData(c *gin.Context, message string, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func GoodResponseWithPage(c *gin.Context, message string, statusCode, total, totalPages, page, Limit int, data interface{}) {
	c.JSON(statusCode, domain.DataPage{
		Status:      true,
		Message:     message,
		Total:       int64(total),
		Pages:       totalPages,
		CurrentPage: uint(page),
		Limit:       uint(Limit),
		Data:        data,
	})
}
