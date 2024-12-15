package service

import (
	"errors"
	"project/domain"
	"project/repository"
	"time"

	"go.uber.org/zap"
)

type ReservationService interface {
	All(timeFilter string) ([]*domain.AllReservation, error)
	Add(reservationRequest *domain.Reservation) error
	GetReservationByID(id uint) (*domain.Reservation, error)
	Update(reservationID uint, updates map[string]interface{}) error
}

type reservationService struct {
	repo repository.ReservationRepository
	log  *zap.Logger
}

func NewReservationService(repo repository.ReservationRepository, log *zap.Logger) ReservationService {
	return &reservationService{repo, log}
}

// All untuk mengambil semua reservasi berdasarkan filter waktu tertentu
func (s *reservationService) All(timeFilter string) ([]*domain.AllReservation, error) {
	reservations, err := s.repo.All(timeFilter)
	if err != nil {
		return nil, err
	}
	if len(reservations) == 0 {
		return nil, errors.New("no reservations found")
	}
	return reservations, nil
}

func (s *reservationService) Add(reservationRequest *domain.Reservation) error {
	// Validasi status hanya boleh Confirmed atau Canceled
	if reservationRequest.Status != "Confirmed" {
		return errors.New("status must be 'Confirmed' ")
	}

	// Validasi Pax Number (maksimal 8 orang)
	if reservationRequest.PaxNumber > 8 {
		return errors.New("pax number cannot exceed 8")
	}

	// Validasi Table Number (maksimal 7 table)
	if reservationRequest.TableNumber > 7 {
		return errors.New("table number cannot exceed 7")
	}

	// Validasi Reservation Date & Time (tidak boleh masa lalu)
	if reservationRequest.ReservationDate.Before(time.Now()) || (reservationRequest.ReservationDate.Equal(time.Now()) && reservationRequest.ReservationTime.Before(time.Now().Local().Truncate(time.Minute))) {
		return errors.New("reservation date and time cannot be in the past")
	}

	// Memanggil fungsi repository untuk menambah reservasi
	err := s.repo.Add(reservationRequest)
	if err != nil {
		s.log.Error("Failed to add reservation", zap.Error(err))
		return err
	}

	s.log.Info("Reservation added successfully")
	return nil
}

func (s *reservationService) GetReservationByID(id uint) (*domain.Reservation, error) {
	reservation, err := s.repo.GetByID(id)
	if err != nil {
		s.log.Error("Failed to fetch reservation by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	return reservation, nil
}

func (s *reservationService) Update(reservationID uint, updates map[string]interface{}) error {
	// Validasi input updates
	if len(updates) == 0 {
		return errors.New("table number cannot 0")
	}

	// Validasi data sebelum mengirim ke repository
	if status, ok := updates["status"]; ok {
		if status != "Canceled" {
			return errors.New("status can only be updated to 'Canceled'")
		}
	}

	if tableNumber, ok := updates["table_number"]; ok {
		if tableNum, valid := tableNumber.(int); valid && tableNum > 7 {
			return errors.New("table number cannot exceed 7")
		}
	}

	// Panggil repository untuk update data
	err := s.repo.Update(reservationID, updates)
	if err != nil {
		s.log.Error("Failed to update reservation", zap.Uint("id", reservationID), zap.Error(err))
		return err
	}

	s.log.Info("Reservation updated successfully", zap.Uint("id", reservationID))
	return nil
}
