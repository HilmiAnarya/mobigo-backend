package agreement

import (
	"context"
	"mobigo-backend/internal/domain"
)

type Repository interface {
	CreateAgreement(ctx context.Context, agreement *domain.Agreement) error
}
