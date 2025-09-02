package schedule

import (
	"context"
	"gorm.io/gorm"
	"mobigo-backend/internal/domain"
)

type gormRepository struct {
	db *gorm.DB
}

func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateSchedule(ctx context.Context, schedule *domain.Schedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}
