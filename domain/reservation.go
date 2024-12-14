package domain

import (
	"time"
)

type Reservation struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ReservationDate time.Time `gorm:"not null" json:"reservation_date" example:"2024-12-14"`
	ReservationTime time.Time `gorm:"not null" json:"reservation_time" example:"14:00:00"`
	TableNumber     uint      `gorm:"not null" json:"table_number"`
	Status          string    `gorm:"not null" json:"status"`
	ReservationName string    `gorm:"size:100;not null" json:"reservation_name"`
	PaxNumber       uint      `gorm:"not null" json:"pax_number"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
}

// ReservationSeed untuk menambahkan contoh data reservasi
func ReservationSeed() []Reservation {
	return []Reservation{
		{
			ReservationDate: time.Date(2024, 12, 14, 0, 0, 0, 0, time.UTC),
			ReservationTime: time.Date(2024, 12, 14, 14, 0, 0, 0, time.UTC),
			TableNumber:     5,
			Status:          "Confirmed",
			ReservationName: "John Doe",
			PaxNumber:       4,
		},
		{
			ReservationDate: time.Date(2024, 12, 15, 0, 0, 0, 0, time.UTC),
			ReservationTime: time.Date(2024, 12, 15, 19, 30, 0, 0, time.UTC),
			TableNumber:     3,
			Status:          "Canceled",
			ReservationName: "Alice Smith",
			PaxNumber:       2,
		},
	}
}
