package payment

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	// For now, this is an admin-only endpoint. We need to protect it.
	// We will add the middleware later once the JWT is fully implemented for staff.
	router.HandleFunc("/api/payments/generate-plan", h.generatePlanHandler).Methods("POST")
}

type generatePlanRequest struct {
	AgreementID        int64   `json:"agreement_id"`
	TotalPrice         float64 `json:"total_price"`
	DownPayment        float64 `json:"down_payment"`
	Tenor              int     `json:"tenor"`
	AnnualInterestRate float64 `json:"annual_interest_rate"`
}

func (h *Handler) generatePlanHandler(w http.ResponseWriter, r *http.Request) {
	var req generatePlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	serviceReq := GeneratePlanRequest{
		AgreementID:        req.AgreementID,
		TotalPrice:         req.TotalPrice,
		DownPayment:        req.DownPayment,
		Tenor:              req.Tenor,
		AnnualInterestRate: req.AnnualInterestRate,
	}

	if err := h.service.GenerateInstallmentPlan(r.Context(), serviceReq); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Installment plan generated successfully"})
}
