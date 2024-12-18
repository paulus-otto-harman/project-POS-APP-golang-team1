package handler

import (
	"mime/multipart"
	"net/http"
	"project/domain"
	"project/helper"
	"project/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
	service service.Service
	logger  *zap.Logger
}

func NewUserController(service service.Service, logger *zap.Logger) *UserController {
	return &UserController{service: service, logger: logger}
}

// Check Email endpoint
// @Summary Check Email
// @Description email must be valid when users want to reset their passwords
// @Tags Auth
// @Accept  json
// @Produce  json
// @Success 200 {object} handler.Response "email is valid"
// @Failure 404 {object} handler.Response "user not found"
// @Router  /users [get]
func (ctrl *UserController) All(c *gin.Context) {
	var err error
	queryURL := c.Query

	// Parse 'page' query parameter with default value
	pageQuery := queryURL("page")
	page := uint(1) // Default value
	if pageQuery != "" {
		page, err = helper.Uint(pageQuery)
		if err != nil {
			ctrl.logger.Error("Invalid 'page' query parameter", zap.Error(err))
			BadResponse(c, "Invalid 'page' query parameter", http.StatusBadRequest)
			return
		}
		if page < 1 {
			page = 1
		}
	}

	// Parse 'limit' query parameter with default value
	limitQuery := queryURL("limit")
	limit := uint(10) // Default value
	if limitQuery != "" {
		limit, err = helper.Uint(limitQuery)
		if err != nil {
			ctrl.logger.Error("Invalid 'limit' query parameter", zap.Error(err))
			BadResponse(c, "Invalid 'limit' query parameter", http.StatusBadRequest)
			return
		}
		if limit < 1 {
			limit = 10
		}
	}

	// Parse sorting parameters
	sortField := queryURL("sort_by")
	sortDirection := queryURL("sort")

	// Fetch users
	users, count, err := ctrl.service.User.All(sortField, sortDirection, page, limit)
	if err != nil {
		ctrl.logger.Error("Failed to fetch users", zap.Error(err))
		BadResponse(c, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Calculate total pages properly
	totalPages := int((count + int64(limit) - 1) / int64(limit)) // Round up division

	// Send success response
	GoodResponseWithPage(c, "Users fetched successfully", http.StatusOK, int(count), totalPages, int(page), int(limit), users)
}

// Registration endpoint
// @Summary Staff Registration
// @Description register staff
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param domain.User body domain.User true " "
// @Success 200 {object} handler.Response "login successfully"
// @Failure 500 {object} handler.Response "server error"
// @Router  /register [post]
func (ctrl *UserController) Registration(c *gin.Context) {
	var file multipart.File
	var fileHeader *multipart.FileHeader
	var filename string
	var err error

	fileHeader, err = c.FormFile("profile_picture")
	if err == nil {

		file, err = fileHeader.Open()
		if err != nil {
			ctrl.logger.Error("Failed to open file", zap.Error(err))
			BadResponse(c, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		filename = fileHeader.Filename
		ctrl.logger.Info("Received new file", zap.String("filename", filename))
	}
	if fileHeader == nil {
		ctrl.logger.Error("File profile photo is missing")
		BadResponse(c, "File profile photo is required", http.StatusBadRequest)
		return
	}
	if err != nil {
		ctrl.logger.Error("Failed to get file from request", zap.Error(err))
		BadResponse(c, "Failed get data: "+err.Error(), http.StatusBadRequest)
		return
	}

	var user domain.User
	birthDateStr := c.PostForm("birth_date")
	if birthDateStr == "" {
		ctrl.logger.Error("Birth date is required")
		BadResponse(c, "Birth date is required", http.StatusBadRequest)
		return
	}
	user.BirthDate = helper.MonthDate(birthDateStr)
	if err := c.ShouldBind(&user); err != nil {
		ctrl.logger.Error("Failed to bind JSON", zap.Error(err))
		BadResponse(c, "invalid request body", http.StatusBadRequest)
		return
	}

	if file != nil {
		newPhotoProfileURL, err := ctrl.service.Category.UploadIcon(file, filename)
		if err != nil {
			BadResponse(c, "Failed to upload new icon: "+err.Error(), http.StatusInternalServerError)
			return
		}
		user.ProfilePhoto = newPhotoProfileURL
	}

	defaultPassword, err := helper.GenerateDefaultPassword(8)
	if err != nil {
		ctrl.logger.Error("Failed to generate default password", zap.Error(err))
		BadResponse(c, "Failed to generate default password", http.StatusInternalServerError)
		return
	}
	if user.Role == "admin" {
		user.Password = defaultPassword
	}

	if err := ctrl.service.User.Register(&user); err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
		TODO :
		- if user is an admin sent email default password
	*/

	if user.Role == "admin" {
		emailData := struct {
			DefaultPassword string
		}{
			DefaultPassword: defaultPassword,
		}
		_, err = ctrl.service.Email.Send(user.Email, "Welcome To COSYPOS, Your Account Have Been Created", "defaultPassword", emailData)
		if err != nil {
			ctrl.logger.Error("Failed to send email", zap.Error(err))
			BadResponse(c, "failed to send email", http.StatusInternalServerError)
			return
		}
	}

	GoodResponseWithData(c, "user registered", http.StatusCreated, user)
}

func (ctrl *UserController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadResponse(c, "invalid parameter", http.StatusBadRequest)
		return
	}

	var newPassword NewPassword
	if err = c.ShouldBindJSON(&newPassword); err != nil {
		BadResponse(c, "invalid password", http.StatusUnprocessableEntity)
		return
	}

	if err = ctrl.service.UpdatePassword(id, newPassword.Password); err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "password changed", http.StatusOK, nil)
}

type NewPassword struct {
	Password        string `binding:"required" json:"password"`
	ConfirmPassword string `binding:"required,eqfield=Password" json:"confirm_password"`
}
