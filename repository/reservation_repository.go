package repository

import (
	"errors"
	"project/domain"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReservationRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewReservationRepository(db *gorm.DB, log *zap.Logger) *ReservationRepository {
	return &ReservationRepository{db: db, log: log}
}

// Create untuk menambahkan reservasi baru
func (repo *ReservationRepository) Add(reservation *domain.Reservation) error {
	// Validasi Status
	if reservation.Status != "Confirmed" && reservation.Status != "Canceled" {
		return errors.New("status must be either 'Confirmed' or 'Canceled'")
	}

	// Validasi Pax Number (max 8)
	if reservation.PaxNumber > 8 {
		return errors.New("pax number cannot exceed 8")
	}

	// Validasi Table Number (max 7)
	if reservation.TableNumber > 7 {
		return errors.New("table number cannot exceed 7")
	}

	// Validasi Reservation Date & Time (tidak boleh masa lalu)
	if reservation.ReservationDate.Before(time.Now()) || (reservation.ReservationDate.Equal(time.Now()) && reservation.ReservationTime.Before(time.Now().Local().Truncate(time.Minute))) {
		return errors.New("reservation date and time cannot be in the past")
	}

	// Validasi Table Number (tidak boleh ada reservasi lain di waktu yang sama)
	var existingReservation domain.Reservation
	err := repo.db.Where("reservation_date = ? AND reservation_time = ? AND table_number = ?", reservation.ReservationDate, reservation.ReservationTime, reservation.TableNumber).
		First(&existingReservation).Error
	if err == nil {
		// Table already reserved
		return errors.New("this table is already reserved at the selected time")
	}

	reservation.ReservationName = reservation.FirstName + " " + reservation.Surname

	// Create new reservation
	err = repo.db.Create(&reservation).Error
	if err != nil {
		repo.log.Error("Failed to create reservation", zap.Error(err))
		return err
	}

	repo.log.Info("Reservation Success")
	return nil
}

// Filter berdasarkan waktu
func (repo *ReservationRepository) filterByTimeQuery(query *gorm.DB, filter string) *gorm.DB {
	now := time.Now().UTC() // Pastikan menggunakan UTC untuk konsistensi waktu

	switch filter {
	case "today":
		// Mulai dari awal hari (00:00:00) hingga akhir hari (23:59:59) UTC
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second)
		query = query.Where("reservation_date BETWEEN ? AND ?", startOfDay, endOfDay)
	case "this_week":
		// Mulai dari hari pertama minggu ini (Senin) hingga hari terakhir minggu ini (Minggu)
		// Dimulai dari Senin pada minggu ini hingga Minggu (di UTC)
		startOfWeek := now.AddDate(0, 0, -int(now.Weekday())+1) // Menggunakan Senin sebagai awal minggu
		endOfWeek := startOfWeek.AddDate(0, 0, 6)               // Menggunakan Minggu sebagai akhir minggu
		query = query.Where("reservation_date BETWEEN ? AND ?", startOfWeek, endOfWeek)
	case "this_month":
		// Mulai dari awal bulan hingga akhir bulan (UTC)
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)
		query = query.Where("reservation_date BETWEEN ? AND ?", startOfMonth, endOfMonth)
	case "this_year":
		// Mulai dari awal tahun hingga akhir tahun (UTC)
		startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		endOfYear := time.Date(now.Year(), 12, 31, 23, 59, 59, 999999999, time.UTC)
		query = query.Where("reservation_date BETWEEN ? AND ?", startOfYear, endOfYear)
	default:
		// Jika tidak ada filter yang sesuai, kita tidak menambah filter waktu
	}

	return query
}

// All untuk mengambil semua reservasi berdasarkan waktu tertentu tanpa pagination
func (repo *ReservationRepository) All(timeFilter string) ([]*domain.AllReservation, error) {
	var reservations []*domain.AllReservation

	// Query awal
	query := repo.db.Model(&domain.Reservation{})

	// Terapkan filter waktu
	query = repo.filterByTimeQuery(query, timeFilter)

	// Ambil semua data reservasi yang sudah terurut berdasarkan tanggal reservasi
	err := query.Order("reservation_date ASC").Find(&reservations).Error
	if err != nil {
		repo.log.Error("Failed to fetch reservations", zap.Error(err))
		return nil, err
	}

	// Jika tidak ada data yang ditemukan
	if len(reservations) == 0 {
		repo.log.Warn("No reservations found")
		return nil, errors.New("no reservations found")
	}

	return reservations, nil
}
