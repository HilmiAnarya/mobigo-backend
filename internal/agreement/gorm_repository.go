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
