// File: internal/booking/handler.go

package booking

import (
	"encoding/json"
	"mobigo-backend/internal/domain"
	"mobigo-backend/pkg/middleware"
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
	r.Use(authMiddleware)

	r.HandleFunc("", h.getAllBookingsHandler).Methods("GET")
	r.HandleFunc("", h.createBookingHandler).Methods("POST")
	r.HandleFunc("/{id}", h.getBookingByIDHandler).Methods("GET")
	r.HandleFunc("/{id}/confirm", h.confirmScheduleHandler).Methods("POST")
	r.HandleFunc("/{id}/decline", h.declineBookingHandler).Methods("PUT")
	r.HandleFunc("/{id}/status", h.updateBookingStatusHandler).Methods("PUT")
}

// --- Request/Response Structs ---
type createBookingRequest struct {
	VehicleID    int64  `json:"vehicle_id"`
	ProposedTime string `json:"proposed_time"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

type confirmScheduleRequest struct {
	Notes string `json:"notes"`
}

type declineBookingRequest struct {
	Reason string `json:"reason"`
}

// --- Handlers ---
func (h *Handler) createBookingHandler(w http.ResponseWriter, r *http.Request) {
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
	proposedTime, err := time.Parse(time.RFC3339, req.ProposedTime)
	if err != nil {
		http.Error(w, "Invalid proposed_time format. Please use RFC3339.", http.StatusBadRequest)
		return
	}
	booking, err := h.service.CreateBooking(r.Context(), userID, req.VehicleID, proposedTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func (h *Handler) confirmScheduleHandler(w http.ResponseWriter, r *http.Request) {
	staffID, _ := r.Context().Value(middleware.UserIDKey).(int64)
	vars := mux.Vars(r)
	bookingID, _ := strconv.ParseInt(vars["id"], 10, 64)

	var req confirmScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	schedule, err := h.service.ConfirmSchedule(r.Context(), bookingID, staffID, req.Notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}

func (h *Handler) declineBookingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, _ := strconv.ParseInt(vars["id"], 10, 64)

	var req declineBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedBooking, err := h.service.DeclineBooking(r.Context(), bookingID, req.Reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedBooking)
}

func (h *Handler) updateBookingStatusHandler(w http.ResponseWriter, r *http.Request) {
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
	newStatus := domain.BookingStatus(req.Status)
	if newStatus != domain.BookingStatusCancelled {
		http.Error(w, "This endpoint can only be used to cancel a booking.", http.StatusBadRequest)
		return
	}
	updatedBooking, err := h.service.UpdateBookingStatus(r.Context(), id, newStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedBooking)
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
