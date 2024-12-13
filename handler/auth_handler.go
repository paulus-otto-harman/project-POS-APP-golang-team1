package handler

import (
	"net/http"
	"project/domain"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthController struct {
	service service.AuthService
	logger  *zap.Logger
}

func NewAuthController(service service.AuthService, logger *zap.Logger) *AuthController {
	return &AuthController{service: service, logger: logger}
}

// Login endpoint
// @Summary User login
// @Description authenticate user
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param domain.User body domain.User true " "
// @Success 200 {object} handler.Response "user authenticated"
// @Failure 401 {object} handler.Response "invalid username and/or password"
// @Failure 500 {object} handler.Response "server error"
// @Router  /login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		BadResponse(c, "invalid request body", http.StatusBadRequest)
		return
	}

	token, isAuthenticated, err := ctrl.service.Login(user)
	if !isAuthenticated {
		BadResponse(c, err.Error(), http.StatusUnauthorized)
		return
	}

	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "user authenticated", http.StatusOK, gin.H{"token": token})
}
