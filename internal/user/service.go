// File: internal/user/service.go
// This is a NEW FILE. It contains the core business logic for the user feature.

package user

import (
	"context"
	"errors"
	"mobigo-backend/internal/domain"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Service defines the business logic operations for users.
type Service interface {
	RegisterStaff(ctx context.Context, fullName, email, password, phoneNumber, address string) (*domain.User, error)
}

// service is the implementation of the Service interface.
type service struct {
	userRepo       Repository
	contextTimeout time.Duration
}

// NewService creates a new instance of the user service.
func NewService(repo Repository, timeout time.Duration) Service {
	return &service{
		userRepo:       repo,
		contextTimeout: timeout,
	}
}

// RegisterStaff handles the business logic for creating a new staff user.
func (s *service) RegisterStaff(ctx context.Context, fullName, email, password, phoneNumber, address string) (*domain.User, error) {
	// 1. Check if user with the same email already exists.
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// This is a system error (e.g., database down).
		return nil, err
	}
	if existingUser != nil {
		// This is a business rule violation, not a system error.
		return nil, errors.New("user with this email already exists")
	}

	// 2. Hash the password for security.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Get the 'staff' role from the database.
	staffRole, err := s.userRepo.GetRoleByName(ctx, "staff")
	if err != nil {
		return nil, err
	}
	if staffRole == nil {
		// This indicates a setup problem - the 'staff' role should always exist.
		return nil, errors.New("staff role not found in database")
	}

	// 4. Create the new user domain object.
	newUser := &domain.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hashedPassword),
		PhoneNumber:  phoneNumber,
		Address:      address,
		Roles:        []*domain.Role{staffRole}, // Assign the staff role
	}

	// 5. Save the new user to the database.
	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// The user object now has the ID and timestamps from the database.
	return newUser, nil
}
