package booking

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

// NewHandler creates a new instance of the booking handler.
func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes sets up the routing for the booking feature.
func (h *Handler) RegisterRoutes(router *mux.Router) {
	r := router.PathPrefix("/api/bookings").Subrouter()
	r.HandleFunc("", h.getAllBookingsHandler).Methods("GET")
}

// getAllBookingsHandler handles retrieving all bookings.
func (h *Handler) getAllBookingsHandler(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.service.ListAllBookings(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve bookings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookings)
}
