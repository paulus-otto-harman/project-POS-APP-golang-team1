package domain

import (
	"encoding/json"
	"fmt"
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
	DepositFee      float64   `gorm:"type:decimal(10,2);not null" json:"deposit_fee,omitempty"` // added DepositFee
	Title           string    `gorm:"size:10" json:"title,omitempty"`                           // added Title
	FirstName       string    `gorm:"size:50" json:"first_name,omitempty"`                      // added FirstName
	Surname         string    `gorm:"size:50" json:"surname,omitempty"`                         // added Surname
	PhoneNumber     string    `gorm:"size:20" json:"phone_number,omitempty"`                    // added PhoneNumber
	EmailAddress    string    `gorm:"size:100" json:"email_address,omitempty"`                  // added EmailAddress
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at" swaggerignore:"true"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at" swaggerignore:"true"`
}
type AllReservation struct {
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

func (r *Reservation) UnmarshalJSON(b []byte) error {
	// Custom structure untuk menampung JSON input
	type Alias Reservation
	aux := &struct {
		ReservationDate string `json:"reservation_date"`
		ReservationTime string `json:"reservation_time"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	// Parse reservation_date
	if aux.ReservationDate != "" {
		parsedDate, err := time.Parse("2006-01-02", aux.ReservationDate)
		if err != nil {
			return fmt.Errorf("invalid reservation_date format, expected YYYY-MM-DD")
		}
		r.ReservationDate = parsedDate
	}

	// Parse reservation_time
	if aux.ReservationTime != "" {
		parsedTime, err := time.Parse("15:04:05", aux.ReservationTime)
		if err != nil {
			return fmt.Errorf("invalid reservation_time format, expected HH:mm:ss")
		}
		r.ReservationTime = parsedTime
	}

	return nil
}

// MarshalJSON untuk menyesuaikan format output JSON
func (r *Reservation) MarshalJSON() ([]byte, error) {
	// Custom structure untuk response JSON
	type Alias Reservation
	return json.Marshal(&struct {
		ReservationDate string `json:"reservation_date"`
		ReservationTime string `json:"reservation_time"`
		*Alias
	}{
		Alias: (*Alias)(r),
		// Format reservation_date menjadi YYYY-MM-DD
		ReservationDate: r.ReservationDate.Format("2006-01-02"),
		// Format reservation_time menjadi HH:mm:ss
		ReservationTime: r.ReservationTime.Format("15:04:05"),
	})
}
