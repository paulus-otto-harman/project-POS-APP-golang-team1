package handler

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"project/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

// @Summary Export Sales Data as CSV
// @Description Export the monthly sales data including total orders, sales, and revenue as a downloadable CSV file.
// @Tags Dashboard
// @Accept  json
// @Produce text/csv
// @Success 200 {string} string "CSV file generated successfully"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /dashboard/export-sales-csv [get]
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		return true
	},
}

// @Summary Real-time Sales Data via WebSocket
// @Description Establish a WebSocket connection to receive real-time sales data including monthly sales, total orders, and revenue.
// @Tags Dashboard
// @Accept  json
// @Produce application/json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /dashboard/sales-data-ws [get]
func (h *DashboardController) SalesDataWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	defer conn.Close()

	// Kirim data secara periodik (misal setiap 5 detik)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Ambil data dashboard dari service
		dashboardSummary, err := h.service.GetDashboard()
		if err != nil {
			h.logger.Error("Failed to fetch dashboard data", zap.Error(err))
			conn.WriteJSON(gin.H{"error": "Failed to fetch dashboard data"})
			continue
		}

		// Format data untuk dikirim
		data := map[string]interface{}{
			"month":             "2024-12", // Ganti dengan data bulan sebenarnya
			"monthlyOrderCount": dashboardSummary.MonthlyOrderCount,
			"monthlySales":      dashboardSummary.MonthlySales,
			"revenue":           float64(dashboardSummary.MonthlyOrderCount) * dashboardSummary.MonthlySales,
		}

		// Kirim data ke klien
		if err := conn.WriteJSON(data); err != nil {
			h.logger.Error("Failed to send data to WebSocket client", zap.Error(err))
			break
		}
	}
}
