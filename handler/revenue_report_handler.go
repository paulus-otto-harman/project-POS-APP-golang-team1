package handler

import (
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RevenueHandler struct {
	service service.RevenueService
	logger  *zap.Logger
}

func NewRevenueHandler(service service.RevenueService, logger  *zap.Logger) *RevenueHandler {
	return &RevenueHandler{service: service, logger: logger}
}

func (h *RevenueHandler) GetTotalRevenueByStatus(c *gin.Context) {
	
}

func (h *RevenueHandler) GetMonthlyRevenue(c *gin.Context) {
	
}

func (h *RevenueHandler) GetProductRevenueDetails(c *gin.Context) {
	
}