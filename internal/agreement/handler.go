package agreement

import (
	"encoding/json"
	"mobigo-backend/internal/domain"
	"mobigo-backend/internal/payment" // Import payment to use its service
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service        Service
	paymentService payment.Service // Add payment service as a dependency
}

func NewHandler(s Service, ps payment.Service) *Handler {
	return &Handler{
		service:        s,
		paymentService: ps,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	r := router.PathPrefix("/api/agreements").Subrouter()
	r.Use(authMiddleware)
	r.HandleFunc("", h.createAgreementHandler).Methods("POST")
}

type createAgreementRequest struct {
	BookingID   int64              `json:"booking_id"`
	FinalPrice  float64            `json:"final_price"`
	PaymentType domain.PaymentType `json:"payment_type"`
	Terms       string             `json:"terms"`
}

// This handler now orchestrates calls to two different services.
func (h *Handler) createAgreementHandler(w http.ResponseWriter, r *http.Request) {
	var req createAgreementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Step 1: Call the agreement service to create the agreement.
	agreement, err := h.service.CreateAgreement(r.Context(), req.BookingID, req.FinalPrice, req.PaymentType, req.Terms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: Check the payment type and call the payment service if needed.
	if agreement.PaymentType == domain.PaymentTypeFull {
		if err := h.paymentService.CreateFullPaymentForAgreement(r.Context(), agreement.ID); err != nil {
			// In a real app, we might want to "roll back" the agreement creation if this fails.
			http.Error(w, "Agreement created, but failed to create full payment record.", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agreement)
}
