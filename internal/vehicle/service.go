package vehicle

import (
	"context"
	"mobigo-backend/internal/domain"
	"time"
)

// Service defines the business logic operations for vehicles.
type Service interface {
	CreateVehicle(ctx context.Context, make, model, vin, description, status string, year int, price float64) (*domain.Vehicle, error)
	GetAllVehicles(ctx context.Context) ([]*domain.Vehicle, error)
	GetVehicleByID(ctx context.Context, id int64) (*domain.Vehicle, error)
	UpdateVehicle(ctx context.Context, id int64, make, model, vin, description, status string, year int, price float64) (*domain.Vehicle, error)
	DeleteVehicle(ctx context.Context, id int64) error
}

// service is the implementation of the Service interface.
type service struct {
	repo           Repository
	contextTimeout time.Duration
}

// NewService creates a new instance of the vehicle service.
func NewService(repo Repository, timeout time.Duration) Service {
	return &service{
		repo:           repo,
		contextTimeout: timeout,
	}
}

// CreateVehicle handles the business logic for creating a new vehicle.
func (s *service) CreateVehicle(ctx context.Context, make, model, vin, description, status string, year int, price float64) (*domain.Vehicle, error) {
	newVehicle := &domain.Vehicle{
		Make:        make,
		Model:       model,
		Year:        year,
		VIN:         vin,
		Price:       price,
		Description: description,
		Status:      status,
	}

	err := s.repo.CreateVehicle(ctx, newVehicle)
	if err != nil {
		return nil, err
	}

	return newVehicle, nil
}

// GetAllVehicles retrieves all vehicles.
func (s *service) GetAllVehicles(ctx context.Context) ([]*domain.Vehicle, error) {
	return s.repo.GetAllVehicles(ctx)
}

// GetVehicleByID retrieves a single vehicle by its ID.
func (s *service) GetVehicleByID(ctx context.Context, id int64) (*domain.Vehicle, error) {
	return s.repo.GetVehicleByID(ctx, id)
}

// UpdateVehicle handles the business logic for updating an existing vehicle.
func (s *service) UpdateVehicle(ctx context.Context, id int64, make, model, vin, description, status string, year int, price float64) (*domain.Vehicle, error) {
	// First, get the existing vehicle to make sure it exists.
	vehicleToUpdate, err := s.repo.GetVehicleByID(ctx, id)
	if err != nil {
		return nil, err // Pass through any database errors
	}
	if vehicleToUpdate == nil {
		return nil, nil // Or a custom "not found" error
	}

	// Update the fields
	vehicleToUpdate.Make = make
	vehicleToUpdate.Model = model
	vehicleToUpdate.Year = year
	vehicleToUpdate.VIN = vin
	vehicleToUpdate.Price = price
	vehicleToUpdate.Description = description
	vehicleToUpdate.Status = status

	// Save the updated vehicle back to the database
	err = s.repo.UpdateVehicle(ctx, vehicleToUpdate)
	if err != nil {
		return nil, err
	}

	return vehicleToUpdate, nil
}

// DeleteVehicle handles the business logic for deleting a vehicle.
func (s *service) DeleteVehicle(ctx context.Context, id int64) error {
	// We could add business logic here, e.g., check if the vehicle is in an active booking.
	// For now, we just pass the call to the repository.
	return s.repo.DeleteVehicle(ctx, id)
}
