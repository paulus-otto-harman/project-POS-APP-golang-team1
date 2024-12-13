package handler

import (
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
// @Success 200 {object} handler.PaginatedResponse "fetch success"
// @Failure 404 {object} handler.Response "categories not found"
// @Failure 500 {object} handler.Response "internal server error"
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

// @Summary Create Category
// @Description Create a new category with an icon
// @Tags Categories
// @Accept  multipart/form-data
// @Produce json
// @Param name formData string true "Category name"
// @Param description formData string false "Category description"
// @Param icon formData file true "Category icon"
// @Success 201 {object} handler.Response "create success"
// @Failure 400 {object} handler.Response "Invalid input"
// @Failure 500 {object} handler.Response "Internal server error"
// @Router /categories/ [post]
func (ctrl *CategoryController) Create(c *gin.Context) {

	fileHeader, err := c.FormFile("icon")
	if fileHeader == nil {
		ctrl.logger.Error("File icon is missing")
		BadResponse(c, "File icon is required", http.StatusBadRequest)
		return
	}
	if err != nil {
		ctrl.logger.Error("Failed to get file from request", zap.Error(err))
		BadResponse(c, "Failed get data: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctrl.logger.Error("Failed to open file", zap.Error(err))
		BadResponse(c, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var category domain.Category
	if err := c.ShouldBind(&category); err != nil {
		ctrl.logger.Error("invalid input", zap.Error(err))
		BadResponse(c, "Invalid category data: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = ctrl.service.Create(&category, file, fileHeader.Filename)
	if err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "create success", http.StatusCreated, nil)
}
