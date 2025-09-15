package payment

import (
	"context"
	"errors"
	"fmt"
	"mobigo-backend/internal/agreement"
	"mobigo-backend/internal/booking"
	"mobigo-backend/internal/domain"
	"mobigo-backend/internal/installment"
	"mobigo-backend/internal/vehicle"
	"time"
)

type Service interface {
	GenerateInstallmentPlan(ctx context.Context, req GeneratePlanRequest) error
	InitiatePayment(ctx context.Context, paymentID int64, customerID int64) (*domain.Payment, error)
	// THE FIX: A new service method specifically for creating the full payment.
	CreateFullPaymentForAgreement(ctx context.Context, agreementID int64) error
}

type service struct {
	paymentRepo     Repository
	installmentRepo installment.Repository
	vehicleRepo     vehicle.Repository
	agreementRepo   agreement.Repository
	bookingRepo     booking.Repository
}

func NewService(paymentRepo Repository, installmentRepo installment.Repository, vehicleRepo vehicle.Repository, agreementRepo agreement.Repository, bookingRepo booking.Repository) Service {
	return &service{
		paymentRepo:     paymentRepo,
		installmentRepo: installmentRepo,
		vehicleRepo:     vehicleRepo,
		agreementRepo:   agreementRepo,
		bookingRepo:     bookingRepo,
	}
}

// CreateFullPaymentForAgreement handles creating the payment and updating the vehicle status for a full payment deal.
func (s *service) CreateFullPaymentForAgreement(ctx context.Context, agreementID int64) error {
	agreement, err := s.agreementRepo.GetByID(ctx, agreementID)
	if err != nil || agreement == nil {
		return errors.New("agreement not found")
	}
	if agreement.PaymentType != domain.PaymentTypeFull {
		return errors.New("agreement is not for a full payment")
	}

	// Create a single payment record
	fullPayment := &domain.Payment{
		AgreementID:   agreementID,
		Amount:        agreement.FinalPrice,
		PaymentMethod: "Full Payment",
		Status:        domain.PaymentStatusPending,
	}
	if err := s.paymentRepo.CreatePayment(ctx, fullPayment); err != nil {
		return err
	}

	// Update vehicle status
	booking, _ := s.bookingRepo.GetBookingByID(ctx, agreement.BookingID)
	vehicleToUpdate, _ := s.vehicleRepo.GetVehicleByID(ctx, booking.VehicleID)
	if vehicleToUpdate != nil {
		vehicleToUpdate.Status = domain.VehicleStatusSold
		return s.vehicleRepo.UpdateVehicle(ctx, vehicleToUpdate)
	}

	return nil
}

// ... (GenerateInstallmentPlan and InitiatePayment methods remain the same)
// ...
type GeneratePlanRequest struct {
	AgreementID        int64
	DownPayment        float64
	Tenor              int
	AnnualInterestRate float64
}

func (s *service) GenerateInstallmentPlan(ctx context.Context, req GeneratePlanRequest) error {
	agreement, err := s.agreementRepo.GetByID(ctx, req.AgreementID)
	if err != nil || agreement == nil {
		return errors.New("agreement not found")
	}
	if agreement.PaymentType != domain.PaymentTypeInstallment {
		return errors.New("this agreement is not for an installment plan")
	}
	existingPayments, err := s.paymentRepo.GetPaymentsByAgreementID(ctx, req.AgreementID)
	if err != nil {
		return err
	}
	if len(existingPayments) > 0 {
		return errors.New("an installment plan already exists for this agreement")
	}
	if req.DownPayment >= agreement.FinalPrice {
		return errors.New("down payment must be less than the total price")
	}

	booking, err := s.bookingRepo.GetBookingByID(ctx, agreement.BookingID)
	if err != nil || booking == nil {
		return errors.New("booking not found for this agreement")
	}
	vehicleToUpdate, err := s.vehicleRepo.GetVehicleByID(ctx, booking.VehicleID)
	if err != nil || vehicleToUpdate == nil {
		return errors.New("vehicle not found for this agreement")
	}

	dpPayment := &domain.Payment{
		AgreementID:   req.AgreementID,
		Amount:        req.DownPayment,
		PaymentMethod: "Down Payment",
		Status:        domain.PaymentStatusPending,
	}
	if err := s.paymentRepo.CreatePayment(ctx, dpPayment); err != nil {
		return err
	}

	loanPrincipal := agreement.FinalPrice - req.DownPayment
	totalInterest := loanPrincipal * (req.AnnualInterestRate / 100) * (float64(req.Tenor) / 12)
	totalRepayment := loanPrincipal + totalInterest
	monthlyBill := totalRepayment / float64(req.Tenor)

	installmentPayment := &domain.Payment{
		AgreementID:   req.AgreementID,
		Amount:        totalRepayment,
		PaymentMethod: "Installment",
		Status:        domain.PaymentStatusPending,
	}
	if err := s.paymentRepo.CreatePayment(ctx, installmentPayment); err != nil {
		return err
	}

	var installmentsToCreate []*domain.Installment
	for i := 1; i <= req.Tenor; i++ {
		dueDate := time.Now().AddDate(0, i, 0)
		inst := &domain.Installment{
			PaymentID:     installmentPayment.ID,
			DueDate:       dueDate,
			AmountDue:     monthlyBill,
			PenaltyAmount: 0,
			TotalDue:      monthlyBill,
			Status:        domain.InstallmentStatusPending,
		}
		installmentsToCreate = append(installmentsToCreate, inst)
	}
	if err := s.installmentRepo.CreateInstallments(ctx, installmentsToCreate); err != nil {
		return err
	}

	vehicleToUpdate.Status = domain.VehicleStatusOnInstallment
	return s.vehicleRepo.UpdateVehicle(ctx, vehicleToUpdate)
}

func (s *service) InitiatePayment(ctx context.Context, paymentID int64, customerID int64) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil || payment == nil {
		return nil, errors.New("payment record not found")
	}
	agreement, err := s.agreementRepo.GetByID(ctx, payment.AgreementID)
	if err != nil || agreement == nil {
		return nil, errors.New("agreement not found for this payment")
	}
	booking, err := s.bookingRepo.GetBookingByID(ctx, agreement.BookingID)
	if err != nil || booking == nil {
		return nil, errors.New("booking not found for this agreement")
	}
	if booking.UserID != customerID {
		return nil, errors.New("unauthorized: you do not own this booking")
	}

	simulatedMidtransID := fmt.Sprintf("MOBI-TX-%d-%d", paymentID, time.Now().Unix())
	simulatedPaymentURL := fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v1/transactions/%s", simulatedMidtransID)
	payment.MidtransTransactionID = &simulatedMidtransID
	payment.PaymentURL = simulatedPaymentURL
	payment.Status = domain.PaymentStatusPending
	return payment, s.paymentRepo.Update(ctx, payment)
}
