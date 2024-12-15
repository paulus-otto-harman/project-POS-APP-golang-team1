package handler

import (
	"net/http"
	"project/domain"
	"project/service"
	"strconv"

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

// @Summary Create Reservation
// @Description Add a new reservation to the system
// @Tags Reservations
// @Accept  json
// @Produce json
// @Param reservation body domain.Reservation true "Reservation details"
// @Success 201 {object} handler.Response "Reservation successfully created"
// @Failure 400 {object} handler.Response "Invalid input"
// @Failure 500 {object} handler.Response "Internal server error"
// @Router /reservations/ [post]
func (ctrl *ReservationController) Add(c *gin.Context) {
	var reservationRequest domain.Reservation
	if err := c.ShouldBindJSON(&reservationRequest); err != nil {
		BadResponse(c, "Invalid reservation data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Call the service layer to add the reservation
	err := ctrl.service.Add(&reservationRequest)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a successful response with the reservation data
	GoodResponseWithData(c, "Reservation success", http.StatusCreated, nil)
}

// GetReservationByID endpoint
// @Summary Get Reservation By ID
// @Description Get reservation details by reservation ID
// @Tags Reservation
// @Accept  json
// @Produce  json
// @Param id path int true "Reservation ID"
// @Success 200 {object} handler.Response "Successfully fetched reservation"
// @Failure 404 {object} handler.Response "Reservation not found"
// @Failure 500 {object} handler.Response "Internal Server Error"
// @Router  /reservations/{id} [get]
func (ctrl *ReservationController) GetByID(c *gin.Context) {
	// Ambil ID dari parameter path
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// Jika ID tidak valid, kembalikan status 400 (Bad Request)
		BadResponse(c, "invalid reservation ID", http.StatusBadRequest)
		return
	}

	// Panggil service untuk mendapatkan reservasi berdasarkan ID
	reservation, err := ctrl.service.GetReservationByID(uint(id))
	if err != nil {
		// Jika reservasi tidak ditemukan, kembalikan status 404
		if err.Error() == "reservation not found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		// Jika terjadi error lain, kembalikan status 500
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// Kirimkan response dengan data reservasi
	GoodResponseWithData(c, "fetch success", http.StatusOK, reservation)
}

// UpdateReservation endpoint
// @Summary Update Reservation
// @Description Update reservation table number or change status to 'Canceled'
// @Tags Reservation
// @Accept  json
// @Produce  json
// @Param id path int true "Reservation ID"
// @Param body body domain.Reservation true "Fields to update (table_number or status)"
// @Success 200 {object} handler.Response "Reservation updated successfully"
// @Failure 400 {object} handler.Response "Invalid input or validation error"
// @Failure 404 {object} handler.Response "Reservation not found"
// @Failure 500 {object} handler.Response "Internal Server Error"
// @Router /reservations/{id} [put]
func (ctrl *ReservationController) Update(c *gin.Context) {
	// Ambil ID dari path parameter
	idParam := c.Param("id")
	reservationID, err := strconv.Atoi(idParam)
	if err != nil {
		BadResponse(c, "invalid reservation ID", http.StatusBadRequest)
		return
	}

	// Bind JSON langsung ke struct domain.Reservation
	var request domain.Reservation
	if err := c.ShouldBindJSON(&request); err != nil {
		BadResponse(c, "invalid request body", http.StatusBadRequest)
		return
	}

	// Siapkan map untuk fields yang akan diupdate
	updates := map[string]interface{}{}
	if request.TableNumber > 0 { // Hanya tambahkan jika table_number tidak 0
		updates["table_number"] = request.TableNumber
	}
	if request.Status != "" { // Hanya tambahkan jika status tidak kosong
		updates["status"] = request.Status
	}

	// Validasi: Harus ada field yang diupdate
	if len(updates) == 0 {
		BadResponse(c, "Table Number cannot 0", http.StatusBadRequest)
		return
	}

	// Panggil service untuk update data
	err = ctrl.service.Update(uint(reservationID), updates)
	if err != nil {
		if err.Error() == "reservation not found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "status can only be updated to 'Canceled'" ||
			err.Error() == "table number cannot exceed 7" {
			BadResponse(c, err.Error(), http.StatusBadRequest)
			return
		}
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// Berhasil update
	GoodResponseWithData(c, "reservation updated successfully", http.StatusOK, nil)
}
