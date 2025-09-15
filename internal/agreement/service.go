package agreement

import (
	"context"
	"errors"
	"mobigo-backend/internal/booking"
	"mobigo-backend/internal/domain"
	"time"
)

// THE FIX: The service no longer depends on payment or vehicle repositories.
type Service interface {
	CreateAgreement(ctx context.Context, bookingID int64, finalPrice float64, paymentType domain.PaymentType, terms string) (*domain.Agreement, error)
	GetByID(ctx context.Context, id int64) (*domain.Agreement, error)
}
type service struct {
	repo        Repository
	bookingRepo booking.Repository
}

func NewService(repo Repository, bookingRepo booking.Repository) Service {
	return &service{
		repo:        repo,
		bookingRepo: bookingRepo,
	}
}

// THE FIX: This function's only job is now to create the agreement.
// All other logic has been moved to other services.
func (s *service) CreateAgreement(ctx context.Context, bookingID int64, finalPrice float64, paymentType domain.PaymentType, terms string) (*domain.Agreement, error) {
	// --- Validation ---
	booking, err := s.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil || booking == nil {
		return nil, errors.New("invalid booking ID")
	}
	if booking.Status != domain.BookingStatusConfirmed {
		return nil, errors.New("agreement can only be created for confirmed bookings")
	}

	// --- Create the Agreement Record ---
	newAgreement := &domain.Agreement{
		BookingID:     bookingID,
		FinalPrice:    finalPrice,
		PaymentType:   paymentType,
		Terms:         terms,
		AgreementDate: time.Now(),
	}
	if err := s.repo.CreateAgreement(ctx, newAgreement); err != nil {
		return nil, err
	}

	return newAgreement, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*domain.Agreement, error) {
	return s.repo.GetByID(ctx, id)
}
