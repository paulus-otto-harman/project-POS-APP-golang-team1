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
// @Description  Fetch all notifications filtered by status and user ID
// @Tags         Notifications
// @Param        user_id  path    int     true   "User ID"
// @Param        status   query   string  false  "Notification status filter"
// @Success      200      {object}  []domain.Notification
// @Failure      400      {object}  Response "Invalid user ID"
// @Failure      404      {object}  Response "Notifications not found"
// @Router       /notifications/{user_id} [get]
func (ctrl *NotificationController) All(c *gin.Context) {
	status := c.Query("status")

	user_id := c.GetString("user-id")
	userID, err := helper.Uint(user_id)
	if err != nil {
		ctrl.logger.Error("Invalid notification ID", zap.Error(err))
		BadResponse(c, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	notification, err := ctrl.service.Notification.All(userID, status)
	if err != nil {
		ctrl.logger.Error("failed to get notifications", zap.Error(err))
		BadResponse(c, err.Error(), http.StatusNotFound)
		return
	}

	GoodResponseWithData(c, "Notification fetched", http.StatusOK, notification)
}

// SendNotificationLowStock godoc
// @Summary      Send low stock notification
// @Description  Check inventory stock and send low stock notifications if necessary
// @Tags         Notifications
// @Success      200  {object}  Response
// @Failure      500  {object}  Response
// @Router       /notifications/low-stock [post]
func (ctrl *NotificationController) SendNotificationLowStock() {
	err := ctrl.service.Notification.CreateNotificationLowStock()
	if err != nil {
		ctrl.logger.Error("failed to send low stock notification", zap.Error(err))
		return
	}
}

// Update godoc
// @Summary      Update notification status
// @Description  Update the status of a single notification
// @Tags         Notifications
// @Param        id      path      int              true  "Notification ID"
// @Param        status  body      domain.UpdateRequest  true  "Status to update"
// @Success      200     {object}  Response
// @Failure      400     {object}  Response
// @Failure      500     {object}  Response
// @Router       /notifications/{id} [put]
func (ctrl *NotificationController) Update(c *gin.Context) {
	// Parse notification ID from path parameters
	id := c.Param("id")
	notifID, err := helper.Uint(id)
	if err != nil {
		ctrl.logger.Error("Invalid notification ID", zap.Error(err))
		BadResponse(c, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	var req domain.UpdateRequest
	// Bind and validate the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		ctrl.logger.Error("Invalid request body", zap.Error(err))
		BadResponse(c, "Status is required", http.StatusBadRequest)
		return
	}

	// Call the service to update the notification
	err = ctrl.service.Notification.Update(notifID, req.Status)
	if err != nil {
		ctrl.logger.Error("Failed to update notification status", zap.Error(err))
		BadResponse(c, "Failed to update notification", http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "Notification status updated successfully", http.StatusOK, nil)
}

// Delete godoc
// @Summary      Delete notification
// @Description  Delete a notification by ID
// @Tags         Notifications
// @Param        id  path  int  true  "Notification ID"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /notifications/{id} [delete]
func (ctrl *NotificationController) Delete(c *gin.Context) {
	id := c.Param("id")
	notifID, err := helper.Uint(id)
	if err != nil {
		ctrl.logger.Error("failed to parse notification id", zap.Error(err))
		BadResponse(c, err.Error(), http.StatusBadRequest)
	}
	err = ctrl.service.Notification.Delete(notifID)
	if err != nil {
		ctrl.logger.Error("failed to delete notification", zap.Error(err))
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	GoodResponseWithData(c, "Notification deleted", http.StatusOK, nil)
}

// BatchUpdate godoc
// @Summary      Batch update notification statuses
// @Description  Update statuses for multiple notifications
// @Tags         Notifications
// @Param        body  body  domain.BatchUpdateNotifRequest  true  "Notification IDs and new status"
// @Success      200  {object}  Response
// @Failure      400  {object}  Response
// @Failure      500  {object}  Response
// @Router       /notifications/batch [put]
func (ctrl *NotificationController) BatchUpdate(c *gin.Context) {
	var req domain.BatchUpdateNotifRequest
	// Parse the JSON request body
	if err := c.BindJSON(&req); err != nil {
		ctrl.logger.Error("failed to bind request body", zap.Error(err))
		BadResponse(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the input
	if len(req.NotificationIDs) == 0 || req.Status == "" {
		BadResponse(c, "Notification IDs and status are required", http.StatusBadRequest)
		return
	}

	// Call the service to perform the batch update
	err := ctrl.service.Notification.BatchUpdate(req.NotificationIDs, req.Status)
	if err != nil {
		ctrl.logger.Error("failed to batch update notification statuses", zap.Error(err))
		BadResponse(c, "Failed to update notifications", http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "Notifications updated successfully", http.StatusOK, nil)
}
