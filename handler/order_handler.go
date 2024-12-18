package handler

import (
	"net/http"
	"project/domain"
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

type orderRequest struct {
	Name       string             `json:"name" binding:"required"`
	TableID    uint               `json:"table_id" binding:"required"`
	OrderItems []domain.OrderItem `json:"order_items" binding:"required,dive"`
}

// @Summary Create Order
// @Description Create a new order with specified table ID and order items
// @Tags Orders
// @Accept  json
// @Produce json
// @Param input body orderRequest true "Order Input"
// @Success 201 {object} Response "Order created successfully"
// @Failure 400 {object} Response "Invalid input"
// @Failure 500 {object} Response "Internal server error"
// @Router /orders/ [post]
func (ctrl *OrderController) Create(c *gin.Context) {

	var input orderRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		BadResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	err := ctrl.service.CreateOrder(input.Name, input.TableID, input.OrderItems)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "Order created successfully", http.StatusCreated, nil)
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
// @Router /orders/{id} [put]
func (ctrl *OrderController) Update(c *gin.Context) {

	GoodResponseWithData(c, "update success", http.StatusOK, nil)
}

type orderResponse struct {
	OrderID           int                 `json:"order_id" example:"1"`
	Name              string              `json:"name" example:"John Doe"`
	TableName         string              `json:"table_name" example:"Table A"`
	PaymentMethodName string              `json:"payment_method_name" example:"cash"`
	OrderItems        []orderItemResponse `json:"order_items"`
	Total             float64             `json:"total" example:"31.99"`
}

type orderItemResponse struct {
	OrderItemID  uint    `json:"order_item_id" example:"1"`
	Quantity     int     `json:"quantity" example:"2"`
	ProductName  string  `json:"product_name" example:"Potato"`
	ProductPrice float64 `json:"product_price" example:"6.99"`
	SubTotal     float64 `json:"sub_total" example:"31.99"`
}

// @Summary Get All Orders
// @Description Retrieve a list of orders with pagination and filters
// @Tags Orders
// @Accept json
// @Produce json
// @Param page query int false "Page number, default is 1"
// @Param limit query int false "Number of items per page, default is 10"
// @Param name query string false "Filter by customer name"
// @Param code_order query string false "Filter by order code"
// @Param status query string false "Filter by order status"
// @Success 200 {object} domain.DataPage{data=[]orderResponse} "fetch success"
// @Failure 404 {object} Response "Orders not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /orders [get]
func (ctrl *OrderController) AllOrders(c *gin.Context) {
	page, _ := helper.Uint(c.DefaultQuery("page", "1"))
	limit, _ := helper.Uint(c.DefaultQuery("limit", "10"))
	name := c.Query("name")
	codeOrder := c.Query("code_order")
	status := c.Query("status")

	orders, totalItems, err := ctrl.service.AllOrders(int(page), int(limit), name, codeOrder, status)
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

