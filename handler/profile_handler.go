package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"project/database"
	"project/domain"
	"project/helper"
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

// Logout endpoint
// @Summary User logout
// @Description sign out user
// @Description <br/><strong>PLEASE ADD "Bearer " TO AUTHORIZATION HEADER !</strong>
// @Tags Profile
// @Produce  json
// @Success 200 {object} Response "user logged out"
// @Failure 401 {object} Response "invalid authorization header"
// @Failure 401 {object} Response "fail to validate signature or session expired"
// @Router  /logout [post]
// @Security Bearer
func (ctrl *ProfileController) Logout(c *gin.Context) {
	userID := c.GetString("user-id")
	ctrl.cacher.HDel(fmt.Sprintf("user:%d", userID), "role")
	ctrl.cacher.Delete(fmt.Sprintf("user:%s:permission", userID))
	ctrl.logger.Info("User logged out successfully")
	GoodResponseWithData(c, "user logged out", http.StatusOK, nil)
}

// Update endpoint
// @Summary Update User Profile
// @Description update user profile
// @Description <br/><strong>PLEASE ADD "Bearer " TO AUTHORIZATION HEADER !</strong>
// @Tags Profile
// @Accept mpfd
// @Produce  json
// @Param full_name formData string true "update user full name"
// @Param email formData string true "update user email"
// @Param address formData string true "update user address"
// @Param password formData string false "change password"
// @Success 200 {object} Response "user updated"
// @Failure 401 {object} Response "invalid authorization header"
// @Failure 401 {object} Response "Invalid input"
// @Router  /profile [put]
// @Security Bearer
func (ctrl *ProfileController) Update(c *gin.Context) {
	userID, err := helper.Uint(c.GetString("user-id"))
	if err != nil {
		ctrl.logger.Error("Unable to retrieve user ID", zap.Error(err))
		BadResponse(c, err.Error(), http.StatusUnauthorized)
		return
	}

	var profile domain.Profile
	if err = c.ShouldBind(&profile); err != nil {
		ctrl.logger.Error("Invalid input", zap.Error(err))
		BadResponse(c, helper.FormatValidationError(err), http.StatusUnprocessableEntity)
		return
	}

	user := domain.User{
		ID:       userID,
		FullName: profile.FullName,
		Email:    profile.Email,
		Address:  profile.Address,
		Password: profile.Password,
	}

	if err = ctrl.service.User.Update(user); err != nil {
		ctrl.logger.Error("Unable to update user", zap.Error(err))
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "user updated", http.StatusOK, nil)
}
