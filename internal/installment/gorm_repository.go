package installment

import (
	"context"
	"mobigo-backend/internal/domain"
	"time"

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

// FindOverdueInstallments finds installments where the due_date is before today and status is pending or overdue.
func (r *gormRepository) FindOverdueInstallments(ctx context.Context) ([]*domain.Installment, error) {
	var installments []*domain.Installment
	today := time.Now().Truncate(24 * time.Hour) // Get the date at the beginning of the day

	err := r.db.WithContext(ctx).
		Where("due_date < ? AND status IN (?, ?)", today, domain.InstallmentStatusPending, domain.InstallmentStatusOverdue).
		Find(&installments).Error

	if err != nil {
		return nil, err
	}
	return installments, nil
}

// UpdateInstallment saves the changes to an installment record.
func (r *gormRepository) UpdateInstallment(ctx context.Context, installment *domain.Installment) error {
	return r.db.WithContext(ctx).Save(installment).Error
}
