package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"project/helper"
	"project/service"
)

type UserPermissionController struct {
	service service.UserPermissionService
	logger  *zap.Logger
}

func NewUserPermissionController(service service.UserPermissionService, logger *zap.Logger) *UserPermissionController {
	return &UserPermissionController{service: service, logger: logger}
}

func (ctrl *UserPermissionController) Update(c *gin.Context) {
	userID, _ := helper.Uint(c.Param("id"))
	var updatedPermission Permission

	if err := c.ShouldBindJSON(&updatedPermission); err != nil {
		BadResponse(c, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := ctrl.service.Update(userID, updatedPermission.Permissions); err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	GoodResponseWithData(c, "permissions updated", http.StatusOK, nil)
}

type Permission struct {
	Permissions []string `json:"permissions" binding:"required"`
}
