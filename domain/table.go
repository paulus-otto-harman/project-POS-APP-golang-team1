package domain

import (
	"time"
)

type Table struct {
	ID        uint      `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	Name      string    `gorm:"size:10;unique" json:"name"`
	Status    bool      `gorm:"type:boolean;default:true" json:"-" example:"true"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-" swaggerignore:"true"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-" swaggerignore:"true"`
}
