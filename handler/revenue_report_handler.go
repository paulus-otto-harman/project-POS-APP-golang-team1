package handler

import (
	"net/http"
	"project/helper"
	"project/service"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RevenueController struct {
	service service.RevenueService
	logger  *zap.Logger
}

func NewRevenueController(service service.RevenueService, logger *zap.Logger) *RevenueController {
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

// GetMonthlyRevenue handles the monthly revenue chart request
func (h *RevenueController) GetMonthlyRevenue(c *gin.Context) {
	// Retrieve query parameters
	statusPayment := c.Query("status_payment")
	year, err := helper.Uint(c.Query("year"))
	if err != nil {
		h.logger.Error("Invalid year", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	// Set the default year if it's not provided
	if year == 0 {
		year = uint(time.Now().Year())
	}

	// Fetch revenue data for the provided status or all statuses combined
	result, err := h.service.GetMonthlyRevenue(statusPayment, int(year))
	if err != nil {
		h.logger.Error("Failed to get monthly revenue",
			zap.String("status", statusPayment),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch revenue data"})
		return
	}

	// Define the ordered list of months
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}

	// Prepare the formatted data for response
	formattedData := MonthlyRevenueResponse{
		Months:  months,
		Revenue: result,
		Year:    year,
	}

	GoodResponseWithData(c, "Monthly revenue data fetched successfully", http.StatusOK, formattedData)
}

type MonthlyRevenueResponse struct {
	Months  []string           `json:"months"`
	Revenue map[string]float64 `json:"revenue"`
	Year    uint               `json:"year"`
}

func (h *RevenueController) GetProductRevenueDetails(c *gin.Context) {

}
