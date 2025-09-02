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
	RegisterStaff(ctx context.Context, fullName, email, password, phoneNumber, address string) (*domain.User, error)
	RegisterCustomer(ctx context.Context, fullName, email, password, phoneNumber string) (*domain.User, error)
	// We now have two distinct login methods.
	LoginStaff(ctx context.Context, email, password string) (string, error)
	LoginCustomer(ctx context.Context, email, password string) (string, error)
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

// RegisterCustomer handles the business logic for creating a new customer user.
func (s *service) RegisterCustomer(ctx context.Context, fullName, email, password, phoneNumber string) (*domain.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	customerRole, err := s.userRepo.GetRoleByName(ctx, "customer")
	if err != nil {
		return nil, err
	}
	if customerRole == nil {
		return nil, errors.New("customer role not found in database")
	}

	newUser := &domain.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hashedPassword),
		PhoneNumber:  phoneNumber,
		Roles:        []*domain.Role{customerRole},
	}

	err = s.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
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

// LoginStaff authenticates a user AND authorizes them as a staff member.
func (s *service) LoginStaff(ctx context.Context, email, password string) (string, error) {
	// Step 1: Authenticate (check email and password)
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err // System error
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Step 2: Authorize (check if the user has a staff or admin role)
	isStaff := false
	for _, role := range user.Roles {
		if role.Name == "staff" || role.Name == "admin" {
			isStaff = true
			break
		}
	}

	if !isStaff {
		// This is a valid user, but they are not staff. Deny access.
		return "", errors.New("access denied: user is not a staff member")
	}

	// Step 3: Generate JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   user.Roles, // Include roles in the token
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// LoginCustomer authenticates a user as a customer.
func (s *service) LoginCustomer(ctx context.Context, email, password string) (string, error) {
	// Step 1: Authenticate
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Step 2: Authorize (ensure they are a customer)
	isCustomer := false
	for _, role := range user.Roles {
		if role.Name == "customer" {
			isCustomer = true
			break
		}
	}
	if !isCustomer {
		return "", errors.New("access denied: staff cannot log in through customer portal")
	}

	// Step 3: Generate JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Customers can have longer sessions
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
