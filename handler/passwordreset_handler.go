package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"project/domain"
	"project/service"
)

type PasswordResetController struct {
	service service.PasswordResetService
	logger  *zap.Logger
}

func NewPasswordResetController(service service.PasswordResetService, logger *zap.Logger) *PasswordResetController {
	return &PasswordResetController{service: service, logger: logger}
}

// Request OTP endpoint
// @Summary Validate Email
// @Description request OTP to reset password. Email must be valid
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.Response "password reset link sent"
// @Failure 500 {object} handler.Response "failed to reset password"
// @Router  /password-reset [post]
func (ctrl *PasswordResetController) Create(c *gin.Context) {
	var passwordResetToken domain.PasswordResetToken
	if err := c.ShouldBindJSON(&passwordResetToken); err != nil {
		BadResponse(c, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := ctrl.service.Create(&passwordResetToken); err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "user registered", http.StatusCreated, passwordResetToken)
}
