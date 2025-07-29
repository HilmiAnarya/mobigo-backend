package user

import (
	"context"
	"gorm.io/gorm"
	"mobigo-backend/internal/domain"
)

// gormRepository is the GORM implementation of the user.Repository interface.
type gormRepository struct {
	db *gorm.DB
}

// NewGORMRepository creates a new instance of our user repository.
func NewGORMRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// CreateUser inserts a new user record into the database using GORM.
// GORM's `Create` method handles the SQL INSERT statement.
func (r *gormRepository) CreateUser(ctx context.Context, user *domain.User) error {
	// We also use .WithContext(ctx) to pass the request context to the database driver,
	// which is important for handling timeouts and cancellations.
	return r.db.WithContext(ctx).Create(user).Error
}

// GetUserByEmail retrieves a user by their email address using GORM.
// GORM's `Where` and `First` methods build the SELECT ... WHERE ... LIMIT 1 query.
func (r *gormRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		// If the record is not found, GORM returns a specific error.
		// We check for this error and return `nil` for both user and error,
		// because "not found" is an expected outcome, not a system failure.
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		// For any other database error, we return it.
		return nil, err
	}
	return &user, nil
}

// GetRoleByName retrieves a role by its name using GORM.
func (r *gormRepository) GetRoleByName(ctx context.Context, roleName string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.WithContext(ctx).Where("name = ?", roleName).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Role not found
		}
		return nil, err
	}
	return &role, nil
}
