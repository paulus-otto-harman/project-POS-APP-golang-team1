package service

import (
	"fmt"
	"project/domain"
	"project/repository"

	"go.uber.org/zap"
)

type DashboardService interface {
	GetDashboardSummary() (*domain.Dashboard, error)
}

type dashboardService struct {
	repo repository.DashboardRepository
	log  *zap.Logger
}

func NewDashboardService(repo repository.DashboardRepository, log *zap.Logger) DashboardService {
	return &dashboardService{repo, log}
}

func (s *dashboardService) GetDashboardSummary() (*domain.Dashboard, error) {
	// Call the repository method to fetch the dashboard summary
	summary, err := s.repo.GetDashboardSummary()
	if err != nil {
		s.log.Error("Failed to get dashboard summary", zap.Error(err))
		return nil, fmt.Errorf("could not get dashboard summary: %v", err)
	}
	return summary, nil
}
