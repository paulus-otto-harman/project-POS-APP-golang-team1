package handler

import (
	"net/http"
	"project/helper"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InventoryController struct {
	service service.InventoryService
	logger  *zap.Logger
}

func NewInventoryController(service service.InventoryService, logger *zap.Logger) *InventoryController {
	return &InventoryController{service: service, logger: logger}
}

// @Summary Get All Invetory
// @Description Retrieve a list of inventory with pagination
// @Tags Categories
// @Accept  json
// @Produce json
// @Param page query int false "Page number, default is 1"
// @Param limit query int false "Number of items per page, default is 10"
// @Param category_id query string false "Category ID to filter inventory"
// @Success 200 {object} domain.DataPage{data=[]domain.Product} "fetch success"
// @Failure 404 {object} Response "categories not found"
// @Failure 500 {object} Response "internal server error"
// @Router /inventory/ [get]
func (ctrl *InventoryController) All(c *gin.Context) {
	page, _ := helper.Uint(c.DefaultQuery("page", "1"))
	limit, _ := helper.Uint(c.DefaultQuery("limit", "10"))
	categoryID := c.Query("category_id")

	products, totalItems, err := ctrl.service.All(int(page), int(limit), categoryID)
	if err != nil {
		if err.Error() == "products not found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalItems + int64(limit) - 1) / int64(limit)

	GoodResponseWithPage(c, "fetch success", http.StatusOK, int(totalItems), int(totalPages), int(page), int(limit), products)
}
