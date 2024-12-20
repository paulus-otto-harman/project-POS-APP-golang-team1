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
)

type UserController struct {
	service service.Service
	logger  *zap.Logger
}

func NewUserController(service service.Service, logger *zap.Logger) *UserController {
	return &UserController{service: service, logger: logger}
}

// All retrieves a paginated list of users.
// @Summary      Get all users
// @Description  Retrieve a list of users with optional pagination, sorting, and filtering.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page      query     int     false  "Page number (default: 1)"
// @Param        limit     query     int     false  "Number of users per page (default: 10)"
// @Param        sort_by   query     string  false  "Field to sort by"
// @Param        sort      query     string  false  "Sort direction (asc/desc)"
// @Success      200       {object}  Response{data=[]domain.User}
// @Failure      400       {object}  Response{error=string}
// @Failure      500       {object}  Response{error=string}
// @Router       /staffs [get]
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

// Registration registers a new user.
// @Summary      Register a new user
// @Description  Create a new user account with optional profile picture upload.
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Param        profile_picture  formData  file    true   "Profile picture file"
// @Param        birth_date       formData  string  true   "User's birth date"
// @Param        data             body      domain.User  true  "User registration data"
// @Success      201              {object}  Response{data=domain.User}
// @Failure      400              {object}  Response{error=string}
// @Failure      500              {object}  Response{error=string}
// @Router       /staffs [post]
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
			ctrl.logger.Error("Failed to upload file", zap.Error(err))
			BadResponse(c, "Failed to upload new photo profile: "+err.Error(), http.StatusInternalServerError)
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

func (ctrl *UserController) UpdatePassword(c *gin.Context) {
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

	if err = ctrl.service.User.UpdatePassword(id, newPassword.Password); err != nil {
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "password changed", http.StatusOK, nil)
}

type NewPassword struct {
	Password        string `binding:"required" json:"password"`
	ConfirmPassword string `binding:"required,eqfield=Password" json:"confirm_password"`
}

// Delete deletes a user by ID.
// @Summary      Delete a user
// @Description  Remove a user from the system by their ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  Response{message=string}
// @Failure      400  {object}  Response{error=string}
// @Failure      500  {object}  Response{error=string}
// @Router       /staffs/{id} [delete]
func (ctrl *UserController) Delete(c *gin.Context) {
	userID, err := helper.Uint(c.Param("id"))
	if err != nil {
		ctrl.logger.Error("invalid parameter", zap.Error(err))
		BadResponse(c, "Invalid parameter", http.StatusBadRequest)
		return
	}

	err = ctrl.service.User.Delete(userID)
	if err != nil {
		ctrl.logger.Error("Fail to delete user", zap.Error(err))
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	GoodResponseWithData(c, "User deleted", http.StatusOK, nil)
}

// Update updates an existing user.
// @Summary      Update user details
// @Description  Modify user data with optional profile picture upload.
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Param        id               path      int     true   "User ID"
// @Param        profile_picture  formData  file    true   "Profile picture file"
// @Param        birth_date       formData  string  true   "User's birth date"
// @Param        data             body      domain.User  true  "User update data"
// @Success      200              {object}  Response{message=string}
// @Failure      400              {object}  Response{error=string}
// @Failure      500              {object}  Response{error=string}
// @Router       /staffs/{id} [put]
func (ctrl *UserController) Update(c *gin.Context) {
	var userInput domain.User
	var file multipart.File
	var fileHeader *multipart.FileHeader
	var filename string
	var err error

	ID := c.GetString("user-id")
	if ID == "" {
		ID = c.Param("id")
	}

	userID, err := helper.Uint(ID)
	if err != nil {
		ctrl.logger.Error("invalid parameter", zap.Error(err))
		BadResponse(c, "Invalid parameter", http.StatusBadRequest)
		return
	}
	userInput.ID = userID

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

	birthDateStr := c.PostForm("birth_date")
	if birthDateStr == "" {
		ctrl.logger.Error("Birth date is required")
		BadResponse(c, "Birth date is required", http.StatusBadRequest)
		return
	}
	userInput.BirthDate = helper.MonthDate(birthDateStr)
	if err := c.ShouldBind(&userInput); err != nil {
		ctrl.logger.Error("Failed to bind JSON", zap.Error(err))
		BadResponse(c, "invalid request body", http.StatusBadRequest)
		return
	}

	if file != nil {
		newPhotoProfileURL, err := ctrl.service.Category.UploadIcon(file, filename)
		if err != nil {
			ctrl.logger.Error("Failed to upload new photo profile", zap.Error(err))
			BadResponse(c, "Failed to upload new photo profile: "+err.Error(), http.StatusInternalServerError)
			return
		}
		userInput.ProfilePhoto = newPhotoProfileURL
	}

	err = ctrl.service.User.Update(userInput)
	if err != nil {
		ctrl.logger.Error("Fail to update user", zap.Error(err))
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	GoodResponseWithData(c, "User updated", http.StatusOK, nil)
}

// GetByID retrieves a single user by their ID.
// @Summary      Get user by ID
// @Description  Retrieve user details by their unique identifier.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  Response{data=domain.User}
// @Failure      400  {object}  Response{error=string}
// @Failure      500  {object}  Response{error=string}
// @Router       /staffs/{id} [get]
func (ctrl *UserController) GetByID(c *gin.Context) {
	userID, err := helper.Uint(c.Param("id"))
	if err != nil {
		ctrl.logger.Error("invalid parameter", zap.Error(err))
		BadResponse(c, "Invalid parameter", http.StatusBadRequest)
		return
	}

	userInput := domain.User{
		ID: userID,
	}

	user, err := ctrl.service.User.GetByID(userInput)
	if err != nil {
		ctrl.logger.Error("Fail to get user", zap.Error(err))
		BadResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}
	GoodResponseWithData(c, "User retrieved", http.StatusOK, user)
}

func (ctrl *UserController) UpdateShiftSchedule() {
	err := ctrl.service.User.UpdateShift()
	if err != nil {
		ctrl.logger.Error("Fail to update user shift schedule", zap.Error(err))
		return
	}
}
