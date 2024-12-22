package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
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

// Request OTP endpoint
// @Summary Check Email
// @Description request OTP to reset password. Email must be valid
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.Response "OTP sent"
// @Failure 500 {object} handler.Response "failed to send OTP"
// @Router  /otp [post]
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

	passwordResetToken := domain.PasswordResetToken{UserID: user.ID, Email: user.Email, Otp: otp}
	if err = ctrl.service.PasswordReset.Create(&passwordResetToken); err != nil {
		BadResponse(c, "fail to create OTP", http.StatusInternalServerError)
		return
	}

	emailData := EmailData{ID: passwordResetToken.ID, OTP: passwordResetToken.Otp}
	_, err = ctrl.service.Email.Send(passwordResetToken.Email, "Request OTP", "otp", emailData)
	if err != nil {
		ctrl.logger.Error("failed to send email", zap.Error(err))
		BadResponse(c, "failed to send email", http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "OTP sent", http.StatusCreated, nil)
}

type EmailData struct {
	ID  uuid.UUID
	OTP string
}

// Validate OTP endpoint
// @Summary Validasi OTP
// @Description validate an OTP to reset password. If OTP is valid, user may change their password.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.Response "OTP sent"
// @Failure 500 {object} handler.Response "failed to send OTP"
// @Router  /otp/:id [put]
func (ctrl *PasswordResetController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadResponse(c, "invalid parameter", http.StatusBadRequest)
		return
	}

	var resetToken ResetToken
	if err = c.ShouldBindJSON(&resetToken); err != nil {
		BadResponse(c, "invalid OTP", http.StatusUnprocessableEntity)
		return
	}

	if err = ctrl.service.PasswordReset.Validate(id, resetToken.OTP); err != nil {
		BadResponse(c, "failed to validate OTP", http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "OTP is valid", http.StatusOK, gin.H{"id": id.String()})
}

type ResetToken struct {
	OTP string `binding:"required" json:"otp"`
}
