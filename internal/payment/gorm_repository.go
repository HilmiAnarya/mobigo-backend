package payment

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

func (r *gormRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}
