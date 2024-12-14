package service

import (
	"errors"
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type ReservationService interface {
	All(timeFilter string) ([]*domain.Reservation, error)
}

type reservationService struct {
	repo repository.ReservationRepository
	log  *zap.Logger
}

func NewReservationService(repo repository.ReservationRepository, log *zap.Logger) ReservationService {
	return &reservationService{repo, log}
}

// All untuk mengambil semua reservasi berdasarkan filter waktu tertentu
func (s *reservationService) All(timeFilter string) ([]*domain.Reservation, error) {
	reservations, err := s.repo.All(timeFilter)
	if err != nil {
		return nil, err
	}
	if len(reservations) == 0 {
		return nil, errors.New("no reservations found")
	}
	return reservations, nil
}
