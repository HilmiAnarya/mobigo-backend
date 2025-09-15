package agreement

import (
	"context"
	"gorm.io/gorm"
	"mobigo-backend/internal/domain"
)

type gormRepository struct{ db *gorm.DB }

func NewGORMRepository(db *gorm.DB) Repository { return &gormRepository{db: db} }

func (r *gormRepository) CreateAgreement(ctx context.Context, agreement *domain.Agreement) error {
	return r.db.WithContext(ctx).Create(agreement).Error
}

func (r *gormRepository) GetByID(ctx context.Context, id int64) (*domain.Agreement, error) {
	var agreement domain.Agreement
	err := r.db.WithContext(ctx).First(&agreement, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &agreement, nil
}
