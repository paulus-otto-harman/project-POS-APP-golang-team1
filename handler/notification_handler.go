package handler

import (
	"net/http"
	"project/domain"
	"project/helper"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NotificationController struct {
	service service.Service
	logger  *zap.Logger
}

func NewNotificationController(service service.Service, logger *zap.Logger) *NotificationController {
	return &NotificationController{service, logger}
}

// GetNotifications godoc
// @Summary      Get notifications
// @Description  Fetch all notifications filtered by status
// @Tags         Notifications
// @Param        status  query  string  false  "Notification status filter"
// @Success      200     {object}  []domain.Notification
// @Failure      404     {object}  Response
// @Router       /notifications [get]
func (ctrl *NotificationController) GetNotifications(c *gin.Context) {
	status := c.Query("status")
	notification, err := ctrl.service.Notification.GetAllNotifications(status)
	if err != nil {
		ctrl.logger.Error("failed to get notifications", zap.Error(err))
		BadResponse(c, err.Error(), http.StatusNotFound)
		return
	}

	GoodResponseWithData(c, "Notification fetched", http.StatusOK, notification)
}
