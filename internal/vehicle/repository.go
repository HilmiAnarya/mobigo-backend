package vehicle

import (
	"context"
	"mobigo-backend/internal/domain"
)

// Repository is the interface that provides vehicle storage methods.
type Repository interface {
	CreateVehicle(ctx context.Context, vehicle *domain.Vehicle) error
	GetAllVehicles(ctx context.Context) ([]*domain.Vehicle, error)
	GetVehicleByID(ctx context.Context, id int64) (*domain.Vehicle, error)
	UpdateVehicle(ctx context.Context, vehicle *domain.Vehicle) error
	DeleteVehicle(ctx context.Context, id int64) error
}
