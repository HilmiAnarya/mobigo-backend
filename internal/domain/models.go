// File: internal/domain/models.go
// This file is updated to add the ProposedDatetime to the Booking struct.

package domain

import (
	"gorm.io/gorm"
	"time"
)

// --- Status Enums ---
type VehicleStatus string

const (
	VehicleStatusAvailable     VehicleStatus = "available"
	VehicleStatusBooked        VehicleStatus = "booked"
	VehicleStatusSold          VehicleStatus = "sold"
	VehicleStatusOnInstallment VehicleStatus = "on_installment" // THE NEW STATUS
)

type BookingStatus string

const (
	BookingStatusPending             BookingStatus = "pending"
	BookingStatusConfirmed           BookingStatus = "confirmed"
	BookingStatusCancelled           BookingStatus = "cancelled"
	BookingStatusCompleted           BookingStatus = "completed"
	BookingStatusRescheduleRequested BookingStatus = "reschedule_requested" // THE NEW STATUS
)

type ScheduleStatus string

const (
	ScheduleStatusScheduled ScheduleStatus = "scheduled"
	ScheduleStatusCompleted ScheduleStatus = "completed"
	ScheduleStatusNoShow    ScheduleStatus = "no-show"
	ScheduleStatusCancelled ScheduleStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusSettlement PaymentStatus = "settlement"
	PaymentStatusExpire     PaymentStatus = "expire"
	PaymentStatusFailure    PaymentStatus = "failure"
	PaymentStatusCancel     PaymentStatus = "cancel"
)

type PaymentType string

const (
	PaymentTypeFull        PaymentType = "full_payment"
	PaymentTypeInstallment PaymentType = "installment"
)

type InstallmentStatus string

const (
	InstallmentStatusPending InstallmentStatus = "pending"
	InstallmentStatusPaid    InstallmentStatus = "paid"
	InstallmentStatusOverdue InstallmentStatus = "overdue"
	InstallmentStatusFailed  InstallmentStatus = "failed"
)

// --- Main Models ---

type User struct {
	ID             int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName       string           `gorm:"not null" json:"full_name"`
	Email          string           `gorm:"unique;not null" json:"email"`
	PasswordHash   string           `gorm:"not null" json:"-"`
	PhoneNumber    string           `json:"phone_number"`
	Address        string           `json:"address"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	DeletedAt      gorm.DeletedAt   `gorm:"index" json:"-"`
	Roles          []*Role          `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	PaymentMethods []*PaymentMethod `gorm:"foreignKey:UserID" json:"payment_methods,omitempty"`
}

type Role struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"unique;not null" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type PaymentMethod struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64          `gorm:"not null" json:"user_id"`
	MidtransToken string         `gorm:"unique;not null" json:"-"`
	CardType      string         `json:"card_type"`
	MaskedCard    string         `gorm:"not null" json:"masked_card"`
	IsDefault     bool           `gorm:"default:false" json:"is_default"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Vehicle struct {
	ID          int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Make        string          `gorm:"not null" json:"make"`
	Model       string          `gorm:"not null" json:"model"`
	Year        int             `gorm:"not null" json:"year"`
	VIN         string          `gorm:"unique;not null" json:"vin"`
	Price       float64         `gorm:"type:decimal(15,2);not null" json:"price"`
	Description string          `json:"description"`
	Status      VehicleStatus   `gorm:"type:varchar(50);not null;default:'available'" json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
	Images      []*VehicleImage `gorm:"foreignKey:VehicleID" json:"images,omitempty"`
}

type VehicleImage struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	VehicleID int64          `gorm:"not null" json:"vehicle_id"`
	ImageURL  string         `gorm:"not null" json:"image_url"`
	IsPrimary bool           `gorm:"default:false" json:"is_primary"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Booking struct {
	ID               int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID           int64          `gorm:"not null" json:"user_id"`
	VehicleID        int64          `gorm:"not null" json:"vehicle_id"`
	Status           BookingStatus  `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	ProposedDatetime *time.Time     `json:"proposed_datetime,omitempty"`
	DeclineReason    *string        `json:"decline_reason,omitempty"` // THE NEW FIELD
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	User             *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Vehicle          *Vehicle       `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Agreement        *Agreement     `gorm:"foreignKey:BookingID" json:"agreement,omitempty"`
}

type Schedule struct {
	ID                  int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID           int64          `gorm:"unique;not null" json:"booking_id"`
	UserID              int64          `gorm:"not null" json:"user_id"`
	AppointmentDatetime time.Time      `gorm:"not null" json:"appointment_datetime"`
	Notes               string         `json:"notes"`
	Status              ScheduleStatus `gorm:"type:varchar(50);not null;default:'scheduled'" json:"status"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
	User                *User          `gorm:"foreignKey:UserID" json:"staff_member,omitempty"`
}

type Agreement struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	BookingID     int64          `gorm:"unique;not null" json:"booking_id"`
	AgreementDate time.Time      `gorm:"not null" json:"agreement_date"`
	FinalPrice    float64        `gorm:"type:decimal(15,2);not null" json:"final_price"` // CORRECT TYPE: float64
	PaymentType   PaymentType    `gorm:"type:enum('full_payment', 'installment');not null" json:"payment_type"`
	Terms         string         `json:"terms"`
	SignedByUser  bool           `gorm:"default:false" json:"signed_by_user"`
	SignedByStaff bool           `gorm:"default:false" json:"signed_by_staff"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Payment struct {
	ID                    int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	AgreementID           int64          `gorm:"not null" json:"agreement_id"`
	Amount                float64        `gorm:"type:decimal(15,2);not null" json:"amount"` // CORRECT TYPE: float64
	PaymentMethod         string         `gorm:"not null" json:"payment_method"`
	Status                PaymentStatus  `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	MidtransTransactionID *string        `gorm:"unique" json:"midtrans_transaction_id,omitempty"`
	PaymentURL            string         `json:"payment_url,omitempty"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

type Installment struct {
	ID            int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	PaymentID     int64             `gorm:"not null" json:"payment_id"`
	DueDate       time.Time         `gorm:"type:date;not null" json:"due_date"`
	AmountDue     float64           `gorm:"type:decimal(15,2);not null" json:"amount_due"`               // CORRECT TYPE: float64
	PenaltyAmount float64           `gorm:"type:decimal(15,2);not null;default:0" json:"penalty_amount"` // CORRECT TYPE: float64
	TotalDue      float64           `gorm:"type:decimal(15,2);not null" json:"total_due"`                // CORRECT TYPE: float64
	Status        InstallmentStatus `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	PaidDate      *time.Time        `gorm:"type:date" json:"paid_date,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `gorm:"index" json:"-"`
}
