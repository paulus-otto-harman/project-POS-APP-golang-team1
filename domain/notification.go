package domain

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID      uint   `json:"id" gorm:"primaryKey" example:"1"`
	Title   string `json:"title" binding:"required" example:"Low Inventory Alert"`
	Content string `json:"content" binding:"required" example:"This is to notify you that the following items are running low in stock:"`
	// Status    string         `json:"status" gorm:"type:VARCHAR(10);check:status IN ('read', 'unread');default:'unread'" binding:"required" example:"unread"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at" example:"2024-12-01T10:00:00Z"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" example:"2024-12-02T10:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`

	// UserNotifications []UserNotification `gorm:"foreignKey:NotificationID" json:"user_notifications"` // Reference the join table
}

type UserNotification struct {
	UserID         uint           `gorm:"primaryKey" json:"user_id"`         // Composite Primary Key
	NotificationID uint           `gorm:"primaryKey" json:"notification_id"` // Composite Primary Key
	Status         string         `gorm:"type:VARCHAR(10);check:status IN ('read', 'unread');default:'unread'" json:"status" binding:"required" example:"unread"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"` // Ensure consistency
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`                              // Supports soft delete for user notifications
}

func (u *UserNotification) BeforeCreate(db *gorm.DB) error {
	// Ensure that the Status field is valid
	validStatuses := map[string]bool{
		"read":   true,
		"unread": true,
	}
	if !validStatuses[u.Status] {
		return fmt.Errorf("invalid status: %s, must be 'read' or 'unread'", u.Status)
	}

	// Set the CreatedAt field to the current timestamp if not already set
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	// Additional custom logic (if needed)
	// Example: Check if UserID and NotificationID combination already exists
	existingRecord := UserNotification{}
	if err := db.Where("user_id = ? AND notification_id = ?", u.UserID, u.NotificationID).First(&existingRecord).Error; err == nil {
		return fmt.Errorf("duplicate user_id and notification_id combination")
	}

	return nil
}

type BatchUpdateNotifRequest struct {
	NotificationIDs []uint `json:"notification_ids"`
	Status          string `json:"status"`
}

// Define a request struct for validation
type UpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

func NotificationSeed() []Notification {
	return []Notification{
		{
			Title:   "Low Inventory Alert",
			Content: "This is to notify you that the following items are running low in stock:",
			// Status:  "unread",
		},
		{
			Title:   "Low Inventory Alert",
			Content: "This is to notify you that the following items are running low in stock:",
			// Status:  "read",
		},
		{
			Title:   "Low Inventory Alert",
			Content: "This is to notify you that the following items are running low in stock:",
			// Status:  "unread",
		},
	}
}
