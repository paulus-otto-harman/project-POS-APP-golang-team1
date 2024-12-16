package routes

import (
	"project/infra"

	"github.com/gin-gonic/gin"
)

func notificationRoutes(ctx infra.ServiceContext, r *gin.Engine) {
	notifHandler := ctx.Ctl.NotificationHandler
	notifGroup := r.Group("/notifications")

	notifGroup.GET("/", notifHandler.GetNotifications)
	notifGroup.PUT("/:id", notifHandler.UpdateNotificationStatus)
	notifGroup.PUT("/batch", notifHandler.BatchUpdateNotificationStatus)
	notifGroup.DELETE("/:id", notifHandler.DeleteNotification)
	notifGroup.POST("/low-stock", notifHandler.SendNotificationLowStock)
}
