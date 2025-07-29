package domain

import (
	"gorm.io/gorm"
	"time"
)

// User corresponds to the 'users' table
type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName     string         `gorm:"not null" json:"full_name"`
	Email        string         `gorm:"unique;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	PhoneNumber  string         `json:"phone_number"`
	Address      string         `json:"address"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // This enables soft delete
	Roles        []*Role        `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// Role corresponds to the 'roles' table
type Role struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // This enables soft delete
}

// Vehicle corresponds to the 'vehicles' table
type Vehicle struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Make        string         `gorm:"not null" json:"make"`
	Model       string         `gorm:"not null" json:"model"`
	Year        int            `gorm:"not null" json:"year"`
	VIN         string         `gorm:"unique;not null" json:"vin"`
	Price       float64        `gorm:"not null" json:"price"`
	Description string         `json:"description"`
	Status      string         `gorm:"not null;default:'available'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // This enables soft delete
}
