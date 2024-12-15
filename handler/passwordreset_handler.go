package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"log"
	"net/http"
	"project/domain"
	"project/service"
)

type PasswordResetController struct {
	service service.Service
	logger  *zap.Logger
}

func NewPasswordResetController(service service.Service, logger *zap.Logger) *PasswordResetController {
	return &PasswordResetController{service: service, logger: logger}
}

// Reset Password endpoint
// @Summary Password Reset
// @Description request password reset
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.Response "OTP sent"
// @Failure 500 {object} handler.Response "failed to send OTP"
// @Router  /password-reset [post]
func (ctrl *PasswordResetController) Create(c *gin.Context) {
	var requestOTP domain.RequestOTP

	if err := c.ShouldBindJSON(&requestOTP); err != nil {
		BadResponse(c, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := ctrl.service.User.Get(domain.User{Email: requestOTP.Email})
	if err != nil {
		BadResponse(c, "user not found", http.StatusNotFound)
		return
	}

	otp := ctrl.service.Otp.Generate()

	passwordResetToken := domain.PasswordResetToken{Email: user.Email, Otp: otp}
	if err = ctrl.service.PasswordReset.Create(&passwordResetToken); err != nil {
		BadResponse(c, "fail to create OTP", http.StatusInternalServerError)
		return
	}

	var emailId string
	emailId, err = ctrl.service.Email.Send(passwordResetToken.Email, "Request OTP", "otp", struct {
		ID  uuid.UUID
		OTP string
	}{ID: passwordResetToken.ID, OTP: passwordResetToken.Otp})
	if err != nil {
		ctrl.logger.Error("failed to send email", zap.Error(err))
		BadResponse(c, "failed to send email", http.StatusInternalServerError)
		return
	}

	log.Println(emailId)

	GoodResponseWithData(c, "OTP sent", http.StatusCreated, nil)
}
