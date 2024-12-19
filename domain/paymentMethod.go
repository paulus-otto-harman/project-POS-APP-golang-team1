package domain

import (
	"time"
)

type PaymentMethod struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;unique" json:"name" example:"Credit Card"`
	Status    bool      `gorm:"type:boolean;default:true" json:"-" example:"true"` // Active or inactive
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-" swaggerignore:"true"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-" swaggerignore:"true"`
}

func PaymentMethodSeed() []PaymentMethod {
	return []PaymentMethod{
		{
			Name: "Cash",
		},
		{
			Name: "Credit Card",
		},
		{
			Name: "E-Wallet",
		},
	}
}
