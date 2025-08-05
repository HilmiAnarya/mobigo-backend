package booking

import (
	"context"
	"mobigo-backend/internal/domain"
)

// Repository is the interface that provides booking storage methods.
type Repository interface {
	GetAllBookings(ctx context.Context) ([]*domain.Booking, error)
	GetBookingByID(ctx context.Context, id int64) (*domain.Booking, error) // New method
	UpdateBooking(ctx context.Context, booking *domain.Booking) error      // New method
	CreateBooking(ctx context.Context, booking *domain.Booking) error
}
