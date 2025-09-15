package vehicleimage

import (
	"context"
	"errors"
	"mobigo-backend/internal/domain"
	"os"
)

type Service interface {
	CreateVehicleImage(ctx context.Context, vehicleID int64, imageURL string, isPrimary bool) (*domain.VehicleImage, error)
	DeleteVehicleImage(ctx context.Context, id int64) error
	SetPrimaryImage(ctx context.Context, id int64) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateVehicleImage(ctx context.Context, vehicleID int64, imageURL string, isPrimary bool) (*domain.VehicleImage, error) {
	image := &domain.VehicleImage{
		VehicleID: vehicleID,
		ImageURL:  imageURL,
		IsPrimary: isPrimary,
	}

	err := s.repo.Create(ctx, image)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (s *service) DeleteVehicleImage(ctx context.Context, id int64) error {
	// First, get the image to find its file path
	image, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if image == nil {
		return errors.New("image not found")
	}

	// Try to delete the physical file from the server
	// The path stored in DB is like "/uploads/filename.jpg", we need "./uploads/filename.jpg"
	filePath := "." + image.ImageURL
	if err := os.Remove(filePath); err != nil {
		// Log the error but don't stop the process. We still want to remove the DB record.
		// In a real app, you'd have a better logging system.
		// log.Printf("Warning: could not delete file %s: %v", filePath, err)
	}

	// Now, soft-delete the record from the database
	return s.repo.Delete(ctx, id)
}

func (s *service) SetPrimaryImage(ctx context.Context, id int64) error {
	image, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if image == nil {
		return errors.New("image not found")
	}

	// In a transaction, first set all other images for this vehicle to is_primary = false
	if err := s.repo.ResetPrimaryImagesForVehicle(ctx, image.VehicleID); err != nil {
		return err
	}

	// Then, set the chosen image to is_primary = true
	image.IsPrimary = true
	return s.repo.Update(ctx, image)
}
