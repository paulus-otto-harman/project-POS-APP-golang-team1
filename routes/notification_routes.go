package routes

import (
	"project/infra"

	"github.com/gin-gonic/gin"
)

func notificationRoutes(ctx infra.ServiceContext, r *gin.Engine) {
	notifHandler := ctx.Ctl.NotificationHandler
	notifGroup := r.Group("/notifications", ctx.Middleware.Authentication())

	notifGroup.GET("/", notifHandler.GetNotifications)
}
