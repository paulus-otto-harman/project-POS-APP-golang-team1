package service

import (
	"errors"
	"project/domain"
	"project/helper"
	"project/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService interface {
	All(sortField, sortDirection string, page, limit uint) ([]domain.User, int64, error)
	Get(user domain.User) (*domain.User, error)
	Register(user *domain.User) error
	UpdatePassword(id uuid.UUID, newPassword string) error
	Delete(id uint) error
	Update(user domain.User) error
	GetByID(userInput domain.User) (*domain.User, error)
	UpdateShift() error
}

type userService struct {
	repo repository.Repository
	log  *zap.Logger
}

func NewUserService(repo repository.Repository, log *zap.Logger) UserService {
	return &userService{repo, log}
}

func (s *userService) All(sortField, sortDirection string, page, limit uint) ([]domain.User, int64, error) {
	return s.repo.User.All(sortField, sortDirection, page, limit)
}

func (s *userService) Register(user *domain.User) error {
	existedUser, err := s.repo.User.Get(*user)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existedUser != nil && existedUser.ID != 0 {
		return errors.New("user email already exists")
	}

	if user.Role == "admin" {
		user.Password = helper.HashPassword(user.Password)
	}

	err = s.repo.User.Create(user)
	if err != nil {
		s.log.Error("Error creating user", zap.Error(err))
		return err
	}

	if user.Role == "admin" {
		adminPermission := domain.UserPermission{
			UserID:       user.ID,
			PermissionID: 1,
		}

		err = s.repo.UserPermission.Create(adminPermission)
		if err != nil {
			s.log.Error("Error creating admin permission", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *userService) Get(user domain.User) (*domain.User, error) {
	return s.repo.User.Get(user)
}

func (s *userService) UpdatePassword(id uuid.UUID, newPassword string) error {
	passwordResetToken := domain.PasswordResetToken{ID: id}
	if err := s.repo.PasswordReset.Get(&passwordResetToken); err != nil {
		return err
	}

	if passwordResetToken.ValidatedAt == nil {
		return errors.New("password reset token is invalid")
	}

	if passwordResetToken.PasswordResetAt != nil {
		return errors.New("password reset token has expired")
	}

	passwordResetToken.User.Password = helper.HashPassword(newPassword)
	if err := s.repo.User.Update(&passwordResetToken.User); err != nil {
		return err
	}

	passwordResetToken.PasswordResetAt = helper.Ptr(time.Now())
	if err := s.repo.PasswordReset.Update(&passwordResetToken); err != nil {
		return err
	}
	return nil
}

func (s *userService) Delete(id uint) error {
	return s.repo.User.Delete(id)
}

func (s *userService) Update(updatedUser domain.User) error {
	existedUser, err := s.repo.User.Get(domain.User{ID: updatedUser.ID})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}

	if err != nil {
		return err
	}

	mergeExistingUserWithUpdatedUser(existedUser, updatedUser)

	if err = s.repo.User.Update(existedUser); err != nil {
		return err
	}

	return nil
}

func mergeExistingUserWithUpdatedUser(existingUser *domain.User, updatedUser domain.User) {
	existingUser.FullName = shouldUpdate(existingUser.FullName, updatedUser.FullName)
	existingUser.Email = shouldUpdate(existingUser.Email, updatedUser.Email)
	existingUser.Role = shouldUpdate(existingUser.Role, updatedUser.Role)
	existingUser.ProfilePhoto = shouldUpdate(existingUser.ProfilePhoto, updatedUser.ProfilePhoto)
	existingUser.PhoneNumber = shouldUpdate(existingUser.PhoneNumber, updatedUser.PhoneNumber)
	existingUser.Salary = shouldUpdate(existingUser.Salary, updatedUser.Salary)
	existingUser.BirthDate = shouldUpdate(existingUser.BirthDate, updatedUser.BirthDate)
	existingUser.ShiftStart = shouldUpdate(existingUser.ShiftStart, updatedUser.ShiftStart)
	existingUser.ShiftEnd = shouldUpdate(existingUser.ShiftEnd, updatedUser.ShiftEnd)
	existingUser.Address = shouldUpdate(existingUser.Address, updatedUser.Address)
	existingUser.AdditionalDetails = shouldUpdate(existingUser.AdditionalDetails, updatedUser.AdditionalDetails)

	if updatedUser.Password == "" {
		return
	}

	hashedUpdatedPassword := helper.HashPassword(updatedUser.Password)
	if !helper.CheckPassword(hashedUpdatedPassword, existingUser.Password) {
		existingUser.Password = hashedUpdatedPassword
	}
}

func shouldUpdate[T comparable](existing, updated T) T {
	if existing != updated {
		return updated
	}
	return existing
}

func (s *userService) GetByID(userInput domain.User) (*domain.User, error) {
	user, err := s.repo.User.Get(userInput)
	if err != nil {
		return nil, err
	}
	user.Permissions = nil
	return user, nil
}

func (s *userService) UpdateShift() error {
	// Define the shift schedule
	shiftSchedule := []Shift{
		{StartTime: "9am", EndTime: "6pm"},
		{StartTime: "2pm", EndTime: "11pm"},
	}

	// Retrieve all users
	users, _, err := s.repo.User.All("", "", 0, 0)
	if err != nil {
		return err
	}

	// Loop through users to update their shifts
	for _, user := range users {
		// Check the current shift of the user and update accordingly
		switch user.ShiftStart {
		case shiftSchedule[0].StartTime:
			user.ShiftStart = shiftSchedule[1].StartTime
			user.ShiftEnd = shiftSchedule[1].EndTime
		case shiftSchedule[1].StartTime:
			user.ShiftStart = shiftSchedule[0].StartTime
			user.ShiftEnd = shiftSchedule[0].EndTime
		default:
			// Default to the first shift schedule if the user's shift doesn't match
			user.ShiftStart = shiftSchedule[0].StartTime
			user.ShiftEnd = shiftSchedule[0].EndTime
		}

		// Update the user in the repository
		err = s.repo.User.Update(&user)
		if err != nil {
			// Log the error and return
			s.log.Error("Error updating user shift", zap.Error(err), zap.Uint("user_id", user.ID))
			return err
		}
	}

	return nil
}

type Shift struct {
	StartTime string
	EndTime   string
}
