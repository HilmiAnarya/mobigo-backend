package vehicleimage

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

func (r *gormRepository) Create(ctx context.Context, image *domain.VehicleImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *gormRepository) GetByID(ctx context.Context, id int64) (*domain.VehicleImage, error) {
	var image domain.VehicleImage
	err := r.db.WithContext(ctx).First(&image, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &image, nil
}

func (r *gormRepository) Update(ctx context.Context, image *domain.VehicleImage) error {
	return r.db.WithContext(ctx).Save(image).Error
}

func (r *gormRepository) Delete(ctx context.Context, id int64) error {
	// GORM's Delete with a struct containing gorm.DeletedAt will perform a soft delete.
	return r.db.WithContext(ctx).Delete(&domain.VehicleImage{}, id).Error
}

// ResetPrimaryImagesForVehicle sets all images for a given vehicle to is_primary = false.
// This is done within a transaction to ensure data consistency.
func (r *gormRepository) ResetPrimaryImagesForVehicle(ctx context.Context, vehicleID int64) error {
	return r.db.WithContext(ctx).Model(&domain.VehicleImage{}).
		Where("vehicle_id = ?", vehicleID).
		Update("is_primary", false).Error
}
