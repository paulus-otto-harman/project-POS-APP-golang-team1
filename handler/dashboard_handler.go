package handler

import (
	"net/http"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardController struct {
	service service.DashboardService
	logger  *zap.Logger
}

func NewDashboardController(service service.DashboardService, logger *zap.Logger) *DashboardController {
	return &DashboardController{service: service, logger: logger}
}

// @Summary Get Dashboard Summary
// @Description Retrieve the summary of the dashboard (daily sales, monthly sales, total tables)
// @Tags Dashboard
// @Accept  json
// @Produce json
// @Success 200 {object} domain.Dashboard "dashboard summary fetched successfully"
// @Failure 500 {object} Response "internal server error"
// @Router /dashboard/summary [get]
func (h *DashboardController) GetDashboard(c *gin.Context) {
	// Call the service method to get the dashboard summary
	summary, err := h.service.GetDashboard()
	if err != nil {
		h.logger.Error("Failed to get dashboard summary", zap.Error(err))
		BadResponse(c, "Failed to get dashboard summary", http.StatusNotFound)
		return
	}
	GoodResponseWithData(c, "Success", http.StatusOK, summary)
}
