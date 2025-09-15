// File: internal/agreement/handler.go

package agreement

import (
	"encoding/json"
	"mobigo-backend/internal/domain"
	"net/http"

	"github.com/gorilla/mux"
)

// THE FIX: The handler is now much simpler. It no longer needs to know about the payment service.
type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
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

// THE FIX: The handler's only job is to translate the request and call its own service.
func (h *Handler) createAgreementHandler(w http.ResponseWriter, r *http.Request) {
	var req createAgreementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	agreement, err := h.service.CreateAgreement(r.Context(), req.BookingID, req.FinalPrice, req.PaymentType, req.Terms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agreement)
}
