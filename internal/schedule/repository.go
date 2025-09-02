package schedule

import (
	"context"
	"mobigo-backend/internal/domain"
)

type Repository interface {
	CreateSchedule(ctx context.Context, schedule *domain.Schedule) error
}
