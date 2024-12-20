package handler

import (
	"mime/multipart"
	"net/http"
	"project/domain"
	"project/helper"
	"project/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CategoryController struct {
	service service.CategoryService
	logger  *zap.Logger
}

func NewCategoryController(service service.CategoryService, logger *zap.Logger) *CategoryController {
	return &CategoryController{service: service, logger: logger}
}

// @Summary Get All Categories
// @Description Retrieve a list of categories with pagination
// @Tags Categories
// @Accept  json
// @Produce json
// @Param page query int false "Page number, default is 1"
// @Param limit query int false "Number of items per page, default is 10"
// @Success 200 {object} domain.DataPage{data=[]domain.Category} "fetch success"
// @Failure 404 {object} Response "categories not found"
// @Failure 500 {object} Response "internal server error"
// @Router /categories/ [get]
func (ctrl *CategoryController) All(c *gin.Context) {
	page, _ := helper.Uint(c.DefaultQuery("page", "1"))
	limit, _ := helper.Uint(c.DefaultQuery("limit", "10"))

	categories, totalItems, err := ctrl.service.All(int(page), int(limit))
	if err != nil {
		if err.Error() == "categories not found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalItems + int64(limit) - 1) / int64(limit)

	GoodResponseWithPage(c, "fetch success", http.StatusOK, int(totalItems), int(totalPages), int(page), int(limit), categories)
}

type CategoryRequest struct {
	// Icon string `json:"icon" binding:"required" form:"icon"`
	Name        string `json:"name" binding:"required,min=3" form:"name"`
	Description string `json:"description" binding:"required,min=20" form:"description"`
}

// @Summary Create Category
// @Description Create a new category with an icon
// @Tags Categories
// @Accept  multipart/form-data
// @Produce json
// @Param name formData string true "Category name"
// @Param description formData string false "Category description"
// @Param icon formData file true "Category icon"
// @Success 201 {object} Response "create success"
// @Failure 400 {object} Response "Invalid input"
// @Failure 500 {object} Response "Internal server error"
// @Router /categories/create [post]
func (ctrl *CategoryController) Create(c *gin.Context) {

	var file multipart.File
	var fileHeader *multipart.FileHeader
	var filename string
	var err error

	fileHeader, err = c.FormFile("icon")
	if err != nil {
		ctrl.logger.Error("File icon is missing", zap.Error(err))
		BadResponse(c, "File icon is required", http.StatusBadRequest)
		return
	}

	file, err = fileHeader.Open()
	if err != nil {
		ctrl.logger.Error("Failed to open file", zap.Error(err))
		BadResponse(c, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename = fileHeader.Filename
	ctrl.logger.Info("Received new file", zap.String("filename", filename))

	var input CategoryRequest
	if err := c.ShouldBind(&input); err != nil {
		ctrl.logger.Error("Invalid input", zap.Error(err))
		BadResponse(c, helper.FormatValidationError(err), http.StatusBadRequest)
		return
	}

	category := &domain.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	if file != nil {
		newIconURL, err := ctrl.service.UploadIcon(file, filename)
		if err != nil {
			BadResponse(c, "Failed to upload new icon: "+err.Error(), http.StatusInternalServerError)
			return
		}
		category.Icon = newIconURL
	}

	if err := ctrl.service.Create(category); err != nil {
		ctrl.logger.Error("Failed to create category", zap.Error(err))
		BadResponse(c, "Failed to create category: "+err.Error(), http.StatusInternalServerError)
		return
	}

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
func (ctrl *CategoryController) Update(c *gin.Context) {

	id := c.Param("id")

	var category domain.Category
	if err := ctrl.service.FindByID(&category, id); err != nil {
		ctrl.logger.Error("Category not found", zap.Error(err))
		BadResponse(c, "Category not found", http.StatusNotFound)
		return
	}

	var file multipart.File
	var fileHeader *multipart.FileHeader
	var filename string
	var err error

	fileHeader, err = c.FormFile("icon")
	if err != nil {
		ctrl.logger.Error("File icon is missing", zap.Error(err))
		BadResponse(c, "File icon is required", http.StatusBadRequest)
		return
	}

	file, err = fileHeader.Open()
	if err != nil {
		ctrl.logger.Error("Failed to open file", zap.Error(err))
		BadResponse(c, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filename = fileHeader.Filename
	ctrl.logger.Info("Received new file", zap.String("filename", filename))

	var input CategoryRequest
	if err := c.ShouldBind(&input); err != nil {
		ctrl.logger.Error("Invalid input", zap.Error(err))
		BadResponse(c, helper.FormatValidationError(err), http.StatusBadRequest)
		return
	}

	if file != nil {
		newIconURL, err := ctrl.service.UploadIcon(file, filename)
		if err != nil {
			BadResponse(c, "Failed to upload new icon: "+err.Error(), http.StatusInternalServerError)
			return
		}
		category.Icon = newIconURL
	}

	category.Name = input.Name
	category.Description = input.Description

	if err := ctrl.service.Update(&category); err != nil {
		ctrl.logger.Error("Failed to update category", zap.Error(err))
		BadResponse(c, "Failed to update category: "+err.Error(), http.StatusBadRequest)
		return
	}

	GoodResponseWithData(c, "update success", http.StatusOK, category)
}

// @Summary Get All Products
// @Description Retrieve a list of products with pagination
// @Tags Categories
// @Accept  json
// @Produce json
// @Param page query int false "Page number, default is 1"
// @Param limit query int false "Number of items per page, default is 10"
// @Param category_id query string false "Category ID to filter products"
// @Success 200 {object} domain.DataPage{data=[]domain.Product} "fetch success"
// @Failure 404 {object} Response "categories not found"
// @Failure 500 {object} Response "internal server error"
// @Router /products/ [get]
func (ctrl *CategoryController) AllProducts(c *gin.Context) {
	page, _ := helper.Uint(c.DefaultQuery("page", "1"))
	limit, _ := helper.Uint(c.DefaultQuery("limit", "10"))
	categoryID := c.Query("category_id")

	products, totalItems, err := ctrl.service.AllProducts(int(page), int(limit), categoryID)
	if err != nil {
		if err.Error() == "products not found" {
			BadResponse(c, err.Error(), http.StatusNotFound)
			return
		}
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := (totalItems + int64(limit) - 1) / int64(limit)

	GoodResponseWithPage(c, "fetch success", http.StatusOK, int(totalItems), int(totalPages), int(page), int(limit), products)
}
