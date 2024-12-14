package handler

import (
	"net/http"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ReservationController struct {
	service service.ReservationService
	logger  *zap.Logger
}

func NewReservationController(service service.ReservationService, logger *zap.Logger) *ReservationController {
	return &ReservationController{service: service, logger: logger}
}

// GetAllReservations endpoint
// @Summary Get All Reservations
// @Description Get a list of reservations based on the time filter (today, this week, this month, this year)
// @Tags Reservation
// @Accept  json
// @Produce  json
// @Param filter query string false "Time Filter (today, this_week, this_month, this_year)" default(this_month)
// @Success 200 {object} handler.Response "Successfully fetched reservations"
// @Failure 500 {object} handler.Response "Internal Server Error"
// @Failure 404 {object} handler.Response "No reservations found"
// @Router  /reservations [get]
func (ctrl *ReservationController) All(c *gin.Context) {
	// Ambil parameter filter dari query string, dengan default "this_month"
	timeFilter := c.DefaultQuery("filter", "this_month")

	// Panggil service untuk mengambil daftar reservasi berdasarkan filter
	reservations, err := ctrl.service.All(timeFilter)
	if err != nil {
		// Jika tidak ada reservasi yang ditemukan, kembalikan error dengan status 404
		if err.Error() == "no reservations found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		// Jika terjadi error lain, kembalikan error 500
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// Kirimkan response dengan data reservasi yang ditemukan
	GoodResponseWithData(c, "fetch success", http.StatusOK, reservations)
}
