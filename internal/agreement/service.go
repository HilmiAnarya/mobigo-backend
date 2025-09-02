package agreement

import (
	"context"
	"mobigo-backend/internal/domain"
	"time"
)

type Service interface {
	CreateAgreement(ctx context.Context, bookingID int64, finalPrice float64, terms string) (*domain.Agreement, error)
}
type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateAgreement(ctx context.Context, bookingID int64, finalPrice float64, terms string) (*domain.Agreement, error) {
	newAgreement := &domain.Agreement{
		BookingID:     bookingID,
		FinalPrice:    finalPrice,
		Terms:         terms,
		AgreementDate: time.Now(),
	}
	err := s.repo.CreateAgreement(ctx, newAgreement)
	return newAgreement, err
}
