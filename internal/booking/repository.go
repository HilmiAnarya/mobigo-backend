package booking

import (
	"context"
	"mobigo-backend/internal/domain"
)

// Repository is the interface that provides booking storage methods.
type Repository interface {
	GetAllBookings(ctx context.Context) ([]*domain.Booking, error)
}
