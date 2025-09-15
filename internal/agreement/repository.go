package agreement

import (
	"context"
	"mobigo-backend/internal/domain"
)

type Repository interface {
	CreateAgreement(ctx context.Context, agreement *domain.Agreement) error
	// New method to fetch an agreement by its ID
	GetByID(ctx context.Context, id int64) (*domain.Agreement, error)
}
