package vehicleimage

import (
	"context"
	"mobigo-backend/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, image *domain.VehicleImage) error
	// New methods
	GetByID(ctx context.Context, id int64) (*domain.VehicleImage, error)
	Update(ctx context.Context, image *domain.VehicleImage) error
	Delete(ctx context.Context, id int64) error
	// New method to handle setting primary image
	ResetPrimaryImagesForVehicle(ctx context.Context, vehicleID int64) error
}
