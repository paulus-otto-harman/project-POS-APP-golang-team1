package domain

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	SuperAdmin UserRole = "super admin"
	Admin      UserRole = "admin"
	Staff      UserRole = "staff"
)

type User struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName          string         `gorm:"size:100;not null" json:"full_name" example:"John Smith" form:"full_name" binding:"required"`
	Email             string         `gorm:"index:,unique,composite:emaildeletedat" json:"email" binding:"required" form:"email"`
	Password          string         `gorm:"not null;default:''" json:"-" example:"password"`
	Role              UserRole       `gorm:"type:varchar(50);not null" json:"role" example:"admin" form:"role"`
	ProfilePhoto      string         `gorm:"size:255" json:"profile_photo" example:"/profile_photo/john_smith.jpg"`
	PhoneNumber       string         `gorm:"size:20" json:"phone_number" example:"+1 (23) 123 4567" form:"phone_number"`
	Salary            float64        `gorm:"type:decimal(10,2)" json:"salary" example:"22000.00" form:"salary"`
	BirthDate         time.Time      `gorm:"type:date" json:"birth_date" time_format:"2006-01-02"`
	ShiftStart        string         `gorm:"type:varchar(20)" json:"shift_start" example:"9am" form:"shift_start"`
	ShiftEnd          string         `gorm:"type:varchar(20)" json:"shift_end" example:"6pm" form:"shift_end"`
	Age               int            `json:"age,omitempty" example:"30" gorm:"-"`
	Address           string         `json:"address" example:"1st Street" gorm:"size:255" form:"address"`
	AdditionalDetails string         `json:"additional_details" form:"additional_details" gorm:"size:255" example:"details"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index:,unique,composite:emaildeletedat" json:"-"`

	Permissions         []Permission         `gorm:"many2many:user_permissions;" json:"permissions"`
	PasswordResetTokens []PasswordResetToken `json:"-"`
	Notifications       []Notification       `gorm:"many2many:user_notifications" json:"user_notifications"` // Reference the join table
}
