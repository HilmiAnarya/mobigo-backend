package booking

import (
	"context"
	"mobigo-backend/internal/domain"
	"time"
)

// Service defines the business logic operations for bookings.
type Service interface {
	ListAllBookings(ctx context.Context) ([]*domain.Booking, error)
}

type service struct {
	repo           Repository
	contextTimeout time.Duration
}

// NewService creates a new instance of the booking service.
func NewService(repo Repository, timeout time.Duration) Service {
	return &service{
		repo:           repo,
		contextTimeout: timeout,
	}
}

// ListAllBookings retrieves all bookings.
func (s *service) ListAllBookings(ctx context.Context) ([]*domain.Booking, error) {
	return s.repo.GetAllBookings(ctx)
}
