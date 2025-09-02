package booking

import (
	"encoding/json"
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
	r.HandleFunc("/{id}", h.getBookingByIDHandler).Methods("GET")
	r.HandleFunc("/{id}/propose-time", h.proposeTimeHandler).Methods("PUT")
	r.HandleFunc("/{id}/confirm", h.confirmScheduleHandler).Methods("POST")
}

// --- New Create Booking Handler ---

type createBookingRequest struct {
	VehicleID   int64  `json:"vehicle_id"`
	BookingDate string `json:"booking_date"` // Expecting "YYYY-MM-DD" format
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

type proposeTimeRequest struct {
	ProposedTime string `json:"proposed_time"`
}

type confirmScheduleRequest struct {
	Notes string `json:"notes"`
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

func (h *Handler) proposeTimeHandler(w http.ResponseWriter, r *http.Request) {
	customerID, _ := r.Context().Value(middleware.UserIDKey).(int64)
	vars := mux.Vars(r)
	bookingID, _ := strconv.ParseInt(vars["id"], 10, 64)

	var req proposeTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	proposedTime, err := time.Parse(time.RFC3339, req.ProposedTime)
	if err != nil {
		http.Error(w, "Invalid time format, use RFC3339", http.StatusBadRequest)
		return
	}

	_, err = h.service.ProposeSchedule(r.Context(), bookingID, customerID, proposedTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
