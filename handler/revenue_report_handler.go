package handler

import (
	"net/http"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RevenueController struct {
	service service.RevenueService
	logger  *zap.Logger
}

func NewRevenueController(service service.RevenueService, logger  *zap.Logger) *RevenueController {
	return &RevenueController{service: service, logger: logger}
}

func (h *RevenueController) GetTotalRevenueByStatus(c *gin.Context) {
	data, err := h.service.GetTotalRevenueByStatus()
	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	GoodResponseWithData(c, "fetch success", http.StatusOK, data)
}

func (h *RevenueController) GetMonthlyRevenue(c *gin.Context) {
	
}

func (h *RevenueController) GetProductRevenueDetails(c *gin.Context) {
	
}