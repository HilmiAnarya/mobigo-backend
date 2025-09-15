// File: internal/booking/service.go

package booking

import (
	"context"
	"errors"
	"mobigo-backend/internal/domain"
	"mobigo-backend/internal/schedule"
	"mobigo-backend/internal/vehicle"
	"time"
)

type Service interface {
	ListAllBookings(ctx context.Context) ([]*domain.Booking, error)
	GetBookingDetails(ctx context.Context, id int64) (*domain.Booking, error)
	CreateBooking(ctx context.Context, userID, vehicleID int64, proposedTime time.Time) (*domain.Booking, error)
	ConfirmSchedule(ctx context.Context, bookingID, staffID int64, notes string) (*domain.Schedule, error)
	DeclineBooking(ctx context.Context, bookingID int64, reason string) (*domain.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID int64, newStatus domain.BookingStatus) (*domain.Booking, error)
}

type service struct {
	bookingRepo  Repository
	scheduleRepo schedule.Repository
	vehicleRepo  vehicle.Repository
}

func NewService(bookingRepo Repository, scheduleRepo schedule.Repository, vehicleRepo vehicle.Repository) Service {
	return &service{
		bookingRepo:  bookingRepo,
		scheduleRepo: scheduleRepo,
		vehicleRepo:  vehicleRepo,
	}
}

func (s *service) CreateBooking(ctx context.Context, userID, vehicleID int64, proposedTime time.Time) (*domain.Booking, error) {
	vehicle, err := s.vehicleRepo.GetVehicleByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("vehicle not found")
	}
	if vehicle.Status != domain.VehicleStatusAvailable {
		return nil, errors.New("vehicle is not available for booking")
	}

	newBooking := &domain.Booking{
		UserID:           userID,
		VehicleID:        vehicleID,
		Status:           domain.BookingStatusPending,
		ProposedDatetime: &proposedTime,
	}
	err = s.bookingRepo.CreateBooking(ctx, newBooking)
	return newBooking, err
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
	if booking.Status != domain.BookingStatusPending {
		return nil, errors.New("only pending bookings can be confirmed")
	}

	vehicle, err := s.vehicleRepo.GetVehicleByID(ctx, booking.VehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("associated vehicle not found")
	}
	if vehicle.Status != domain.VehicleStatusAvailable {
		return nil, errors.New("vehicle is no longer available")
	}

	vehicle.Status = domain.VehicleStatusBooked
	booking.Status = domain.BookingStatusConfirmed

	newSchedule := &domain.Schedule{
		BookingID:           bookingID,
		UserID:              staffID,
		AppointmentDatetime: *booking.ProposedDatetime,
		Notes:               notes,
		Status:              domain.ScheduleStatusScheduled,
	}

	// In a real app, this should be a single transaction.
	if err := s.vehicleRepo.UpdateVehicle(ctx, vehicle); err != nil {
		return nil, err
	}
	if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
		return nil, err
	}
	if err := s.scheduleRepo.CreateSchedule(ctx, newSchedule); err != nil {
		return nil, err
	}

	return newSchedule, nil
}

func (s *service) DeclineBooking(ctx context.Context, bookingID int64, reason string) (*domain.Booking, error) {
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}
	if booking.Status != domain.BookingStatusPending {
		return nil, errors.New("only pending bookings can be declined")
	}

	booking.Status = domain.BookingStatusRescheduleRequested
	booking.ProposedDatetime = nil
	booking.DeclineReason = &reason

	if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *service) UpdateBookingStatus(ctx context.Context, bookingID int64, newStatus domain.BookingStatus) (*domain.Booking, error) {
	if newStatus != domain.BookingStatusCancelled {
		return nil, errors.New("this action is only for cancelling a booking")
	}

	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if booking == nil {
		return nil, errors.New("booking not found")
	}

	if booking.Status == domain.BookingStatusConfirmed {
		vehicle, err := s.vehicleRepo.GetVehicleByID(ctx, booking.VehicleID)
		if err != nil {
			return nil, err
		}
		if vehicle != nil && vehicle.Status == domain.VehicleStatusBooked {
			vehicle.Status = domain.VehicleStatusAvailable
			if err := s.vehicleRepo.UpdateVehicle(ctx, vehicle); err != nil {
				return nil, err
			}
		}
	}

	booking.Status = newStatus
	if err := s.bookingRepo.UpdateBooking(ctx, booking); err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *service) ListAllBookings(ctx context.Context) ([]*domain.Booking, error) {
	return s.bookingRepo.GetAllBookings(ctx)
}

func (s *service) GetBookingDetails(ctx context.Context, id int64) (*domain.Booking, error) {
	return s.bookingRepo.GetBookingByID(ctx, id)
}
