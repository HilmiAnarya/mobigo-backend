package user

import (
	"context"
	"errors"
	"mobigo-backend/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Service defines the business logic operations for users.
type Service interface {
	RegisterStaff(ctx context.Context, fullName, email, password, phoneNumber, addess string) (*domain.User, error)
	// LoginStaff now returns the user and a JWT token string.
	LoginStaff(ctx context.Context, email, password string) (*domain.User, string, error)
}

// service is the implementation of the Service interface.
type service struct {
	userRepo       Repository
	jwtSecret      []byte // Add a field for the JWT secret key
	contextTimeout time.Duration
}

// NewService creates a new instance of the user service.
func NewService(repo Repository, jwtSecret string, timeout time.Duration) Service {
	return &service{
		userRepo:       repo,
		jwtSecret:      []byte(jwtSecret), // Store the secret as a byte slice
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

func (s *service) LoginStaff(ctx context.Context, email, password string) (*domain.User, string, error) {
	// 1. Find the user by their email address.
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// This is a system error (e.g., database down).
		return nil, "", err
	}
	if user == nil {
		// User not found. Return a generic error to prevent email enumeration attacks.
		return nil, "", errors.New("invalid email or password")
	}

	// 2. Compare the provided password with the stored hash.
	// bcrypt.CompareHashAndPassword is a secure function that prevents timing attacks.
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// If the passwords don't match, bcrypt returns an error.
		return nil, "", errors.New("invalid email or password")
	}

	// 3. Login successful. Return the user object.
	// Create the claims for the token.
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create a new token object, specifying signing method and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key to get the complete, signed token string.
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}
