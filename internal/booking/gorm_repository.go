package booking

import (
	"context"
	"gorm.io/gorm"
	"mobigo-backend/internal/domain"
)

type gormRepository struct {
	db *gorm.DB
}

// NewGORMRepository creates a new instance of our booking repository.
func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// GetAllBookings retrieves all booking records.
// It uses GORM's Preload feature to automatically fetch the related
// User and Vehicle data for each booking, which is very powerful.
func (r *gormRepository) GetAllBookings(ctx context.Context) ([]*domain.Booking, error) {
	var bookings []*domain.Booking
	err := r.db.WithContext(ctx).
		Preload("User").          // Load the associated User
		Preload("Vehicle").       // Load the associated Vehicle
		Order("created_at desc"). // Show newest bookings first
		Find(&bookings).Error
	return bookings, err
}
