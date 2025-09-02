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
// CreateUser now uses a transaction to ensure both the user and their role association are created.
func (r *gormRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		if err := tx.Model(user).Association("Roles").Append(user.Roles); err != nil {
			return err
		}
		return nil
	})
}

func (r *gormRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	// THE FIX: Preload("Roles") tells GORM to also fetch the associated roles for this user.
	// Without this, the user.Roles slice will always be empty.
	err := r.db.WithContext(ctx).Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
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
