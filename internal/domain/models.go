package domain

import "time"

// User corresponds to the 'users' table
type User struct {
	ID           int64     `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Omit password hash from JSON responses
	PhoneNumber  string    `json:"phone_number"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Vehicle corresponds to the 'vehicles' table
type Vehicle struct {
	ID          int64     `json:"id"`
	Make        string    `json:"make"`
	Model       string    `json:"model"`
	Year        int       `json:"year"`
	VIN         string    `json:"vin"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
