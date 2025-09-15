package payment

import (
	"context"
	"mobigo-backend/internal/domain"
)

// Repository defines the interface for payment data operations.
type Repository interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	GetPaymentsByAgreementID(ctx context.Context, agreementID int64) ([]*domain.Payment, error)
	GetByID(ctx context.Context, id int64) (*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
}
