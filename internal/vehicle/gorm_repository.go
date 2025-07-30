package vehicle

import (
	"context"
	"gorm.io/gorm"
	"mobigo-backend/internal/domain"
)

// gormRepository is the GORM implementation of the vehicle.Repository interface.
type gormRepository struct {
	db *gorm.DB
}

// NewGORMRepository creates a new instance of our vehicle repository.
func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// CreateVehicle saves a new vehicle record to the database.
func (r *gormRepository) CreateVehicle(ctx context.Context, vehicle *domain.Vehicle) error {
	return r.db.WithContext(ctx).Create(vehicle).Error
}

// GetAllVehicles retrieves all vehicle records from the database.
// Note: In a real app, you would add pagination here.
func (r *gormRepository) GetAllVehicles(ctx context.Context) ([]*domain.Vehicle, error) {
	var vehicles []*domain.Vehicle
	err := r.db.WithContext(ctx).Find(&vehicles).Error
	return vehicles, err
}

// GetVehicleByID retrieves a single vehicle by its ID.
func (r *gormRepository) GetVehicleByID(ctx context.Context, id int64) (*domain.Vehicle, error) {
	var vehicle domain.Vehicle
	err := r.db.WithContext(ctx).First(&vehicle, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found is an expected outcome
		}
		return nil, err
	}
	return &vehicle, nil
}

// UpdateVehicle modifies an existing vehicle record in the database.
func (r *gormRepository) UpdateVehicle(ctx context.Context, vehicle *domain.Vehicle) error {
	return r.db.WithContext(ctx).Save(vehicle).Error
}

// DeleteVehicle soft-deletes a vehicle record from the database.
// Because our domain.Vehicle model has gorm.DeletedAt, GORM will automatically
// perform a soft delete (updating the deleted_at column).
func (r *gormRepository) DeleteVehicle(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&domain.Vehicle{}, id).Error
}
