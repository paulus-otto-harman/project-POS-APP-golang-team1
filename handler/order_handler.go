package handler

import (
	"net/http"
	"project/helper"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrderController struct {
	service service.OrderService
	logger  *zap.Logger
}

func NewOrderController(service service.OrderService, logger *zap.Logger) *OrderController {
	return &OrderController{service: service, logger: logger}
}

// @Summary Get All Tables
// @Description Retrieve a list of tables with pagination
// @Tags Tables
// @Accept  json
// @Produce json
// @Param page query int false "Page number, default is 1"
// @Param limit query int false "Number of items per page, default is 10"
// @Success 200 {object} domain.DataPage{data=[]domain.Table} "fetch success"
// @Failure 404 {object} Response "tables not found"
// @Failure 500 {object} Response "internal server error"
// @Router /tables/ [get]
func (ctrl *OrderController) AllTables(c *gin.Context) {
	page, _ := helper.Uint(c.DefaultQuery("page", "1"))
	limit, _ := helper.Uint(c.DefaultQuery("limit", "10"))

	tables, totalItems, err := ctrl.service.AllTables(int(page), int(limit))
	if err != nil {
		if err.Error() == "tables not found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalItems + int64(limit) - 1) / int64(limit)

	GoodResponseWithPage(c, "fetch success", http.StatusOK, int(totalItems), int(totalPages), int(page), int(limit), tables)
}

// @Summary Get All Payments
// @Description Retrieve a list of payments
// @Tags Payments
// @Accept  json
// @Produce json
// @Success 200 {object} Response{data=[]domain.PaymentMethod} "fetch success"
// @Failure 404 {object} Response "payments not found"
// @Failure 500 {object} Response "internal server error"
// @Router /payments/ [get]
func (ctrl *OrderController) AllPayments(c *gin.Context) {

	payments, err := ctrl.service.AllPayments()
	if err != nil {
		if err.Error() == "payments not found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "fetch success", http.StatusOK, payments)
}

// @Summary Create Category
// @Description Create a new category with an icon
// @Tags Categories
// @Accept  multipart/form-data
// @Produce json
// @Param name formData string true "Category name"
// @Param description formData string false "Category description"
// @Param icon formData file true "Category icon"
// @Success 201 {object} Response{data=domain.Category} "create success"
// @Failure 400 {object} Response "Invalid input"
// @Failure 500 {object} Response "Internal server error"
// @Router /categories/create [post]
func (ctrl *OrderController) Create(c *gin.Context) {

	GoodResponseWithData(c, "create success", http.StatusCreated, nil)
}

// @Summary Update Category
// @Description Update an existing category with an optional new icon. If no new icon is provided, the existing icon will be retained.
// @Tags Categories
// @Accept  multipart/form-data
// @Produce json
// @Param id path string true "Category ID"
// @Param name formData string false "Category name"
// @Param description formData string false "Category description"
// @Param icon formData file false "New category icon"
// @Success 200 {object} Response{data=domain.Category} "update success"
// @Failure 400 {object} Response "invalid input"
// @Failure 400 {object} Response "file icon is missing"
// @Failure 404 {object} Response "category not found"
// @Failure 500 {object} Response "internal server error"
// @Router /categories/{id} [put]
func (ctrl *OrderController) Update(c *gin.Context) {

	GoodResponseWithData(c, "update success", http.StatusOK, nil)
}

// @Summary Get All Products
// @Description Retrieve a list of orders with pagination
// @Tags Orders
// @Accept  json
// @Produce json
// @Param page query int false "Page number, default is 1"
// @Param limit query int false "Number of items per page, default is 10"
// @Param category_id query string false "Order ID to filter products"
// @Success 200 {object} domain.DataPage{data=[]domain.Order} "fetch success"
// @Failure 404 {object} Response "Order not found"
// @Failure 500 {object} Response "internal server error"
// @Router /products/ [get]
func (ctrl *OrderController) AllOrders(c *gin.Context) {
	page, _ := helper.Uint(c.DefaultQuery("page", "1"))
	limit, _ := helper.Uint(c.DefaultQuery("limit", "10"))
	orderID := c.Query("order_id")

	orders, totalItems, err := ctrl.service.AllOrders(int(page), int(limit), orderID)
	if err != nil {
		if err.Error() == "orders not found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalItems + int64(limit) - 1) / int64(limit)

	GoodResponseWithPage(c, "fetch success", http.StatusOK, int(totalItems), int(totalPages), int(page), int(limit), orders)
}
