package handler

import (
	"encoding/csv"
	"fmt"
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

func (h *DashboardController) ExportSalesDataCSV(c *gin.Context) {
	// Fetch sales data per month from existing service
	dashboardSummary, err := h.service.GetDashboard() // Adjust the limit and page as needed
	if err != nil {
		BadResponse(c, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare CSV writer
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=sales_data.csv")
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write CSV Header
	writer.Write([]string{"Bulan", "Jumlah Order", "Sales", "Revenue"})

	// Convert MonthlyOrderCount to float64 before multiplying
	revenue := float64(dashboardSummary.MonthlyOrderCount) * dashboardSummary.MonthlySales

	// Write data to CSV
	writer.Write([]string{
		"2024-12", // Replace with actual month value if available
		fmt.Sprintf("%d", dashboardSummary.MonthlyOrderCount), // Jumlah Order
		fmt.Sprintf("%.2f", dashboardSummary.MonthlySales),    // Sales
		fmt.Sprintf("%.2f", revenue),                          // Revenue (calculated)
	})
}
