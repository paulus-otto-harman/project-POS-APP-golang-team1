package routes

import (
	"project/infra"

	"github.com/gin-gonic/gin"
)

func notificationRoutes(ctx infra.ServiceContext, r *gin.Engine) {
	notifHandler := ctx.Ctl.NotificationHandler
	notifGroup := r.Group("/notifications")

	notifGroup.GET("/:user_id", notifHandler.All)
	notifGroup.PUT("/:id", notifHandler.Update)
	notifGroup.PUT("/batch", notifHandler.BatchUpdate)
	notifGroup.DELETE("/:id", notifHandler.Delete)
	notifGroup.POST("/low-stock", notifHandler.SendNotificationLowStock)
}
