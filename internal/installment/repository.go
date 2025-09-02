package installment

import (
	"context"
	"mobigo-backend/internal/domain"
)

// Repository defines the interface for installment data operations.
type Repository interface {
	// CreateInstallments saves a slice of new installment records to the database.
	CreateInstallments(ctx context.Context, installments []*domain.Installment) error
}
