package installment

import (
	"context"
	"mobigo-backend/internal/domain"

	"gorm.io/gorm"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// CreateInstallments uses a transaction to ensure all installments are created or none are.
func (r *gormRepository) CreateInstallments(ctx context.Context, installments []*domain.Installment) error {
	return r.db.WithContext(ctx).Create(&installments).Error
}
