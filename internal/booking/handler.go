package booking

import (
	"encoding/json"
	"mobigo-backend/internal/domain"
	"mobigo-backend/pkg/middleware" // Import our new middleware package
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	r := router.PathPrefix("/api/bookings").Subrouter()

	// All routes in this handler are protected by the auth middleware.
	r.Use(authMiddleware)

	r.HandleFunc("", h.getAllBookingsHandler).Methods("GET")
	r.HandleFunc("", h.createBookingHandler).Methods("POST")
	r.HandleFunc("/{id}", h.getBookingByIDHandler).Methods("GET")      // New route
	r.HandleFunc("/{id}/status", h.updateStatusHandler).Methods("PUT") // New route
}

// --- New Create Booking Handler ---

type createBookingRequest struct {
	VehicleID   int64  `json:"vehicle_id"`
	BookingDate string `json:"booking_date"` // Expecting "YYYY-MM-DD" format
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) createBookingHandler(w http.ResponseWriter, r *http.Request) {
	// Get the userID from the context, which was added by our middleware.
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Could not retrieve user ID from token", http.StatusInternalServerError)
		return
	}

	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bookingDate, err := time.Parse("2006-01-02", req.BookingDate)
	if err != nil {
		http.Error(w, "Invalid date format. Please use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	booking, err := h.service.CreateBooking(r.Context(), userID, req.VehicleID, bookingDate)
	if err != nil {
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

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

func (h *Handler) getBookingByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	booking, err := h.service.GetBookingDetails(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve booking", http.StatusInternalServerError)
		return
	}
	if booking == nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func (h *Handler) updateStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert the incoming string to our safe enum type
	newStatus := domain.BookingStatus(req.Status)

	updatedBooking, err := h.service.UpdateBookingStatus(r.Context(), id, newStatus)
	if err != nil {
		if err.Error() == "booking not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update booking status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedBooking)
}
