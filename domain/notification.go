package domain

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey" example:"1"`
	Title     string         `json:"title" binding:"required" example:"Low Inventory Alert"`
	Content   string         `json:"content" binding:"required" example:"This is to notify you that the following items are running low in stock:"`
	Status    string         `json:"status" gorm:"type:VARCHAR(10);check:status IN ('read', 'unread');default:'unread'" binding:"required" example:"unread"`
	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at" example:"2024-12-01T10:00:00Z"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" example:"2024-12-02T10:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
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
			Status:  "unread",
		},
		{
			Title:   "Low Inventory Alert",
			Content: "This is to notify you that the following items are running low in stock:",
			Status:  "read",
		},
		{
			Title:   "Low Inventory Alert",
			Content: "This is to notify you that the following items are running low in stock:",
			Status:  "unread",
		},
	}
}
