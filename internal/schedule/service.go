package schedule

import (
	"context"
	"mobigo-backend/internal/domain"
	"time"
)

type Service interface {
	CreateSchedule(ctx context.Context, bookingID, staffUserID int64, apptTime time.Time, notes string) (*domain.Schedule, error)
}
type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateSchedule(ctx context.Context, bookingID, staffUserID int64, apptTime time.Time, notes string) (*domain.Schedule, error) {
	newSchedule := &domain.Schedule{
		BookingID:           bookingID,
		UserID:              staffUserID, // This is the ID of the logged-in staff member
		AppointmentDatetime: apptTime,
		Notes:               notes,
		Status:              domain.ScheduleStatusScheduled,
	}
	err := s.repo.CreateSchedule(ctx, newSchedule)
	return newSchedule, err
}
