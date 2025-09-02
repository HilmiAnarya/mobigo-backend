package booking

import (
	"context"
	"errors"
	"mobigo-backend/internal/domain"
	"mobigo-backend/internal/schedule"
	"time"
)

type Service interface {
	ListAllBookings(ctx context.Context) ([]*domain.Booking, error)
	GetBookingDetails(ctx context.Context, id int64) (*domain.Booking, error)
	CreateBooking(ctx context.Context, userID, vehicleID int64, bookingDate time.Time) (*domain.Booking, error)
	// New method for customers
	ProposeSchedule(ctx context.Context, bookingID, customerID int64, proposedTime time.Time) (*domain.Booking, error)
	// New method for admins
	ConfirmSchedule(ctx context.Context, bookingID, staffID int64, notes string) (*domain.Schedule, error)
}

type service struct {
	bookingRepo    Repository
	scheduleRepo   schedule.Repository // Add dependency on schedule repository
	contextTimeout time.Duration
}

// NewService creates a new instance of the booking service.
func NewService(bookingRepo Repository, scheduleRepo schedule.Repository, timeout time.Duration) Service {
	return &service{
		bookingRepo:    bookingRepo,
		scheduleRepo:   scheduleRepo,
		contextTimeout: timeout,
	}
}

func (s *service) ProposeSchedule(ctx context.Context, bookingID, customerID int64, proposedTime time.Time) (*domain.Booking, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}
	if booking.UserID != customerID {
		return nil, errors.New("unauthorized")
	}

	booking.ProposedDatetime = &proposedTime
	err = s.bookingRepo.UpdateBooking(ctx, booking)
	return booking, err
}

func (s *service) ConfirmSchedule(ctx context.Context, bookingID, staffID int64, notes string) (*domain.Schedule, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}
	if booking.ProposedDatetime == nil {
		return nil, errors.New("customer has not proposed a time for this booking")
	}

	newSchedule := &domain.Schedule{
		BookingID:           bookingID,
		UserID:              staffID,
		AppointmentDatetime: *booking.ProposedDatetime,
		Notes:               notes,
		Status:              domain.ScheduleStatusScheduled,
	}
	if err := s.scheduleRepo.CreateSchedule(ctx, newSchedule); err != nil {
		return nil, err
	}

	booking.Status = domain.BookingStatusConfirmed
	if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
		return nil, err
	}
	return newSchedule, nil
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

	err := s.bookingRepo.CreateBooking(ctx, newBooking)
	if err != nil {
		return nil, err
	}

	return newBooking, nil
}

// ListAllBookings retrieves all bookings.
func (s *service) ListAllBookings(ctx context.Context) ([]*domain.Booking, error) {
	return s.bookingRepo.GetAllBookings(ctx)
}

func (s *service) GetBookingDetails(ctx context.Context, id int64) (*domain.Booking, error) {
	return s.bookingRepo.GetBookingByID(ctx, id)
}
