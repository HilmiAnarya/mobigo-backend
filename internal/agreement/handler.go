package agreement

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	r := router.PathPrefix("/api/agreements").Subrouter()
	r.Use(authMiddleware)
	r.HandleFunc("", h.createAgreementHandler).Methods("POST")
}

type createAgreementRequest struct {
	BookingID  int64   `json:"booking_id"`
	FinalPrice float64 `json:"final_price"`
	Terms      string  `json:"terms"`
}

func (h *Handler) createAgreementHandler(w http.ResponseWriter, r *http.Request) {
	var req createAgreementRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	agreement, err := h.service.CreateAgreement(r.Context(), req.BookingID, req.FinalPrice, req.Terms)
	if err != nil {
		http.Error(w, "Failed to create agreement", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(agreement)
}
