package installment

import (
	"context"
	"mobigo-backend/internal/domain"
)

// Repository defines the interface for installment data operations.
type Repository interface {
	// CreateInstallments saves a slice of new installment records to the database.
	CreateInstallments(ctx context.Context, installments []*domain.Installment) error
	// FindOverdueInstallments retrieves all installments that are past their due date and not yet paid.
	FindOverdueInstallments(ctx context.Context) ([]*domain.Installment, error)
	// UpdateInstallment updates a single installment record in the database.
	UpdateInstallment(ctx context.Context, installment *domain.Installment) error
}
