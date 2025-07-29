package user

import (
	"context"
	"mobigo-backend/internal/domain"
)

// Repository is the interface that provides user storage methods.
type Repository interface {
	// CreateUser saves a new user to the database.
	CreateUser(ctx context.Context, user *domain.User) error
	// GetUserByEmail finds a user by their email address.
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	// GetRoleByName finds a role by its name (e.g., "admin", "staff").
	GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error)
}
