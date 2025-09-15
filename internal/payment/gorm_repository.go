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

func (r *gormRepository) GetPaymentsByAgreementID(ctx context.Context, agreementID int64) ([]*domain.Payment, error) {
	var payments []*domain.Payment
	err := r.db.WithContext(ctx).Where("agreement_id = ?", agreementID).Find(&payments).Error
	return payments, err
}

func (r *gormRepository) GetByID(ctx context.Context, id int64) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.WithContext(ctx).First(&payment, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &payment, nil
}

func (r *gormRepository) Update(ctx context.Context, payment *domain.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}
