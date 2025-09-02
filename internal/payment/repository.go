package payment

import (
	"context"
	"mobigo-backend/internal/domain"
)

// Repository defines the interface for payment data operations.
type Repository interface {
	CreatePayment(ctx context.Context, payment *domain.Payment) error
}
