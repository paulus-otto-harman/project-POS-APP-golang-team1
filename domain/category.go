package domain

import (
	"time"
)

type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	Icon        string    `gorm:"size:255;not null" json:"icon,omitempty" example:"/icon/category.png"`
	Name        string    `gorm:"size:100;unique" json:"name"`
	Description string    `gorm:"type:text" example:"lorem" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
}
