package handler

import (
	"net/http"
	"project/helper"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryController struct {
	service service.CategoryService
	logger  *zap.Logger
}

func NewCategoryController(service service.CategoryService, logger *zap.Logger) *CategoryController {
	return &CategoryController{service: service, logger: logger}
}

// Check Email endpoint
// @Summary Check Email
// @Description email must be valid when users want to reset their passwords
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.Response "email is valid"
// @Failure 404 {object} handler.Response "user not found"
// @Router  /users [get]
func (ctrl *CategoryController) All(c *gin.Context) {
	page, _ := helper.Uint(c.DefaultQuery("page", "1"))
	limit, _ := helper.Uint(c.DefaultQuery("limit", "10"))

	categories, totalItems, err := ctrl.service.All(int(page), int(limit))
	if err != nil {
		if err.Error() == "categories not found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	
	totalPages := (totalItems + int64(limit) - 1) / int64(limit)

	GoodResponseWithPage(c, "fetch success", http.StatusOK, int(totalItems), int(totalPages), int(page), int(limit), categories)
}
