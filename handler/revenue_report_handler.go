package handler

import (
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
	
}

func (h *RevenueController) GetMonthlyRevenue(c *gin.Context) {
	
}

func (h *RevenueController) GetProductRevenueDetails(c *gin.Context) {
	
}