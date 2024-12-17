package handler

import (
	"fmt"
	"net/http"
	"project/domain"
	"project/helper"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InventoryController struct {
	service service.InventoryService
	logger  *zap.Logger
}

func NewInventoryController(service service.InventoryService, logger *zap.Logger) *InventoryController {
	return &InventoryController{service: service, logger: logger}
}

// @Summary Get All Inventory
// @Description Retrieve a list of inventory with filters and pagination
// @Tags Inventory
// @Accept  json
// @Produce json
// @Param page query int false "Page number, default is 1"
// @Param limit query int false "Number of items per page, default is 10"
// @Param product_status query string false "Product status: 'Active' or 'Inactive'"
// @Param category_name query string false "Category name to filter inventory"
// @Param stock query string false "Stock status: 'In Stock' or 'Out Of Stock'"
// @Param quantity query int false "Specific product quantity"
// @Param min_price query float64 false "Minimum price filter"
// @Param max_price query float64 false "Maximum price filter"
// @Success 200 {object} domain.DataPage{data=[]domain.Inventory} "Fetch success"
// @Failure 404 {object} Response "Inventory not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /inventory/ [get]
func (ctrl *InventoryController) All(c *gin.Context) {
	// Ambil query parameter dengan nilai default jika kosong
	page, _ := helper.Uint(c.DefaultQuery("page", "1"))
	limit, _ := helper.Uint(c.DefaultQuery("limit", "10"))
	productStatus := c.Query("product_status") // "Active", "Inactive", atau kosong
	categoryName := c.Query("category_name")   // Nama kategori
	stock := c.Query("stock")                  // "In Stock" atau "Out Of Stock" atau "Low Stock"
	quantity, _ := helper.Uint(c.DefaultQuery("quantity", "0"))
	minPrice, _ := helper.Float(c.DefaultQuery("min_price", "0"))
	maxPrice, _ := helper.Float(c.DefaultQuery("max_price", "0"))

	// Panggil service untuk ambil data
	products, totalItems, err := ctrl.service.All(
		int(page),
		int(limit),
		productStatus,
		categoryName,
		stock,
		int(quantity),
		minPrice,
		maxPrice,
	)

	// Error Handling
	if err != nil {
		if err.Error() == "inventory not found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hitung total halaman
	totalPages := (totalItems + int64(limit) - 1) / int64(limit)

	// Kirim response
	GoodResponseWithPage(c, "fetch success", http.StatusOK, int(totalItems), int(totalPages), int(page), int(limit), products)
}

// @Summary Add Inventory
// @Description Add a new inventory item
// @Tags Inventory
// @Accept  json
// @Produce json
// @Param inventory body InventoryInput true "Inventory data"
// @Success 201 {object} domain.Inventory "Inventory created successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 404 {object} Response "Category not found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /inventory/ [post]
func (ctrl *InventoryController) Add(c *gin.Context) {
	var input struct {
		CategoryName string  `json:"category_name" binding:"required"`
		Name         string  `json:"name" binding:"required"`
		CodeProduct  string  `json:"code_product" binding:"required"`
		Quantity     int     `json:"quantity" binding:"required"`
		Price        float64 `json:"price" binding:"required"`
		Status       string  `json:"status" binding:"required"`
		Image        string  `json:"image"`
	}

	// Bind input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		ctrl.logger.Warn("Invalid input", zap.Error(err))
		BadResponse(c, "invalid input", http.StatusBadRequest)
		return
	}

	// Map input ke Inventory domain model
	inventory := domain.Inventory{
		Name:        input.Name,
		CodeProduct: input.CodeProduct,
		Quantity:    input.Quantity,
		Price:       input.Price,
		Status:      input.Status,
		Image:       input.Image,
	}

	// Panggil service untuk menambahkan inventory
	_, err := ctrl.service.Add(&inventory, input.CategoryName)
	if err != nil {
		if err.Error() == fmt.Sprintf("category name '%s' not found", input.CategoryName) {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	// Kirim response sukses
	GoodResponseWithData(c, "inventory created successfully", http.StatusCreated, nil)
}

// @Summary Update Inventory
// @Description Update inventory details by ID
// @Tags Inventory
// @Accept  json
// @Produce json
// @Param id path int true "Inventory ID"
// @Param inventory body domain.Inventory true "Inventory data"
// @Param category_name query string false "Category name to update category ID"
// @Success 200 {object} domain.Inventory "Inventory updated successfully"
// @Failure 404 {object} Response "Inventory not found"
// @Failure 400 {object} Response "Bad request"
// @Failure 500 {object} Response "Internal server error"
// @Router /inventory/{id} [put]
func (ctrl *InventoryController) Update(c *gin.Context) {
	id := c.Param("id")
	categoryName := c.DefaultQuery("category_name", "")

	var inventory domain.Inventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		ctrl.logger.Error("Invalid request body", zap.Error(err))
		BadResponse(c, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert id to uint
	inventoryID, err := helper.Uint(id)
	if err != nil {
		ctrl.logger.Error("Invalid inventory ID", zap.String("id", id))
		BadResponse(c, "Invalid inventory ID", http.StatusBadRequest)
		return
	}

	// Panggil service untuk update inventory
	_, err = ctrl.service.Update(inventoryID, &inventory, categoryName)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "Inventory updated successfully", http.StatusOK, nil)
}
