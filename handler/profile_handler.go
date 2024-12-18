package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"project/database"
	"project/infra/jwt"
	"project/service"
)

type ProfileController struct {
	service service.Service
	logger  *zap.Logger
	jwt     jwt.JWT
	cacher  database.Cacher
}

func NewProfileController(service service.Service, logger *zap.Logger, cacher database.Cacher, jwt jwt.JWT) *ProfileController {
	return &ProfileController{service, logger, jwt, cacher}
}

func (ctrl *ProfileController) Logout(c *gin.Context) {
	userID := c.GetString("user-id")
	ctrl.cacher.HDel(fmt.Sprintf("user:%d", userID), "role")
	ctrl.cacher.Delete(fmt.Sprintf("user:%s:permission", userID))
	ctrl.logger.Info("User logged out successfully")
	GoodResponseWithData(c, "user logged out", http.StatusOK, nil)
}
