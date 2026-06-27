package services

import (
	"context"

	"go-api/internal/core/ports"
)

type healthService struct {
	repo ports.HealthRepository
}

func NewHealthService(repo ports.HealthRepository) *healthService {
	return &healthService{
		repo: repo,
	}
}

func (s *healthService) CheckHealth(ctx context.Context) error {
	return s.repo.CheckHealth(ctx)
}
