package booking

import (
	"context"
	"errors"
	"mobigo-backend/internal/domain"
	"time"
)

// Service defines the business logic operations for bookings.
type Service interface {
	ListAllBookings(ctx context.Context) ([]*domain.Booking, error)
	GetBookingDetails(ctx context.Context, id int64) (*domain.Booking, error)                                   // New method
	UpdateBookingStatus(ctx context.Context, id int64, newStatus domain.BookingStatus) (*domain.Booking, error) // New method
	CreateBooking(ctx context.Context, userID, vehicleID int64, bookingDate time.Time) (*domain.Booking, error)
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

// CreateBooking handles the business logic for creating a new booking.
func (s *service) CreateBooking(ctx context.Context, userID, vehicleID int64, bookingDate time.Time) (*domain.Booking, error) {
	newBooking := &domain.Booking{
		UserID:      userID,
		VehicleID:   vehicleID,
		BookingDate: bookingDate,
		// CHANGED: We now use our safe, compile-time constant instead of a raw string.
		Status: domain.BookingStatusPending,
	}

	err := s.repo.CreateBooking(ctx, newBooking)
	if err != nil {
		return nil, err
	}

	return newBooking, nil
}

// ListAllBookings retrieves all bookings.
func (s *service) ListAllBookings(ctx context.Context) ([]*domain.Booking, error) {
	return s.repo.GetAllBookings(ctx)
}

func (s *service) GetBookingDetails(ctx context.Context, id int64) (*domain.Booking, error) {
	return s.repo.GetBookingByID(ctx, id)
}

func (s *service) UpdateBookingStatus(ctx context.Context, id int64, newStatus domain.BookingStatus) (*domain.Booking, error) {
	booking, err := s.repo.GetBookingByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}

	// Here you can add business rules, e.g., a "completed" booking cannot be "cancelled".
	// For now, we'll allow any change.
	booking.Status = newStatus

	if err := s.repo.UpdateBooking(ctx, booking); err != nil {
		return nil, err
	}
	return booking, nil
}
