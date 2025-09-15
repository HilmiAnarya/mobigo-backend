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

func (r *gormRepository) CreateBooking(ctx context.Context, booking *domain.Booking) error {
	return r.db.WithContext(ctx).Create(booking).Error
}

// GetAllBookings retrieves all booking records.
// It uses GORM's Preload feature to automatically fetch the related
// User and Vehicle data for each booking, which is very powerful.
func (r *gormRepository) GetAllBookings(ctx context.Context) ([]*domain.Booking, error) {
	var bookings []*domain.Booking
	err := r.db.WithContext(ctx).
		Preload("User"). // Load the associated User
		Preload("Vehicle"). // Load the associated Vehicle
		Preload("Agreement.Payments"). // Preload Payments related to the Agreement
		Order("created_at desc"). // Show newest bookings first
		Find(&bookings).Error
	return bookings, err
}

func (r *gormRepository) GetBookingByID(ctx context.Context, id int64) (*domain.Booking, error) {
	var booking domain.Booking
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Vehicle").
		Preload("Agreement.Payments"). // THE FIX: Also load the associated agreement.
		First(&booking, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found is an expected outcome
		}
		return nil, err
	}
	return &booking, nil
}

func (r *gormRepository) UpdateBooking(ctx context.Context, booking *domain.Booking) error {
	return r.db.WithContext(ctx).Save(booking).Error
}
