package payment

import (
	"encoding/json"
	"mobigo-backend/pkg/middleware"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	r := router.PathPrefix("/api/payments").Subrouter()
	r.Use(authMiddleware)

	r.HandleFunc("/generate-plan", h.generatePlanHandler).Methods("POST")
	r.HandleFunc("/{id}/initiate", h.initiatePaymentHandler).Methods("POST")
}

// generatePlanRequest uses float64 to match the service and domain layers.
type generatePlanRequest struct {
	AgreementID        int64   `json:"agreement_id"`
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

func (h *Handler) initiatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	customerID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Could not retrieve user ID from token", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	paymentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	payment, err := h.service.InitiatePayment(r.Context(), paymentID, customerID)
	if err != nil {
		if err.Error() == "unauthorized: you do not own this booking" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}
