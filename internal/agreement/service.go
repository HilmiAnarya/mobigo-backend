package agreement

import (
	"context"
	"errors"
	"mobigo-backend/internal/booking"
	"mobigo-backend/internal/domain"
	"time"
)

// THE FIX: We define a contract for what we need from the payment service.
// The agreement service does not know or care who fulfills this contract.
type PaymentCreator interface {
	CreateFullPaymentForAgreement(ctx context.Context, agreementID int64) error
}

type Service interface {
	CreateAgreement(ctx context.Context, bookingID int64, finalPrice float64, paymentType domain.PaymentType, terms string) (*domain.Agreement, error)
	GetByID(ctx context.Context, id int64) (*domain.Agreement, error)
}

type service struct {
	repo           Repository
	bookingRepo    booking.Repository
	paymentCreator PaymentCreator // THE FIX: The service now depends on the interface, not a concrete type.
}

// THE FIX: The constructor now accepts any struct that fulfills the PaymentCreator contract.
func NewService(repo Repository, bookingRepo booking.Repository, pc PaymentCreator) Service {
	return &service{
		repo:           repo,
		bookingRepo:    bookingRepo,
		paymentCreator: pc,
	}
}

// THE FIX: The service now contains the full business logic, orchestrated correctly.
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

	// --- LOGIC BRANCH based on Payment Type ---
	if paymentType == domain.PaymentTypeFull {
		// Call the payment creation logic via the interface.
		if err := s.paymentCreator.CreateFullPaymentForAgreement(ctx, newAgreement.ID); err != nil {
			// In a real app, we might want to "roll back" the agreement creation if this fails.
			return nil, errors.New("agreement created, but failed to create full payment record")
		}
	}

	return newAgreement, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*domain.Agreement, error) {
	return s.repo.GetByID(ctx, id)
}
