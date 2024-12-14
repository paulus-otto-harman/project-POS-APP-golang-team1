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
func (repo *ReservationRepository) Create(reservation *domain.Reservation) error {
	if err := repo.db.Create(&reservation).Error; err != nil {
		repo.log.Error("Failed to create reservation", zap.Error(err))
		return err
	}
	return nil
}

// Filter berdasarkan waktu
func (repo *ReservationRepository) filterByTimeQuery(query *gorm.DB, filter string) *gorm.DB {
	now := time.Now()

	switch filter {
	case "today":
		// Mulai dari awal hari (00:00:00) hingga akhir hari (23:59:59)
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second)
		query = query.Where("reservation_date BETWEEN ? AND ?", startOfDay, endOfDay)
	case "this_week":
		// Mulai dari hari pertama minggu ini (Senin) hingga hari terakhir minggu ini (Minggu)
		startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
		endOfWeek := startOfWeek.AddDate(0, 0, 6)
		query = query.Where("reservation_date BETWEEN ? AND ?", startOfWeek, endOfWeek)
	case "this_month":
		// Mulai dari awal bulan hingga akhir bulan
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)
		query = query.Where("reservation_date BETWEEN ? AND ?", startOfMonth, endOfMonth)
	case "this_year":
		// Mulai dari awal tahun hingga akhir tahun
		startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endOfYear := time.Date(now.Year(), 12, 31, 23, 59, 59, 999999999, now.Location())
		query = query.Where("reservation_date BETWEEN ? AND ?", startOfYear, endOfYear)
	default:
		// Jika tidak ada filter yang sesuai, kita tidak menambah filter waktu
	}

	return query
}

// All untuk mengambil semua reservasi berdasarkan waktu tertentu tanpa pagination
func (repo *ReservationRepository) All(timeFilter string) ([]*domain.Reservation, error) {
	var reservations []*domain.Reservation

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
