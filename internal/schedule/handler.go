package schedule

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mobigo-backend/pkg/middleware"
	"net/http"
	"time"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(router *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	r := router.PathPrefix("/api/schedules").Subrouter()
	r.Use(authMiddleware)
	r.HandleFunc("", h.createScheduleHandler).Methods("POST")
}

type createScheduleRequest struct {
	BookingID   int64  `json:"booking_id"`
	StaffUserID int64  `json:"staff_user_id"`
	ApptTime    string `json:"appointment_time"` // e.g., "2025-08-15T14:00:00Z"
	Notes       string `json:"notes"`
}

func (h *Handler) createScheduleHandler(w http.ResponseWriter, r *http.Request) {
	// Get the logged-in staff member's ID from the context.
	staffUserID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		http.Error(w, "Could not retrieve user ID from token", http.StatusInternalServerError)
		return
	}

	var req createScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	apptTime, err := time.Parse(time.RFC3339, req.ApptTime)
	if err != nil {
		http.Error(w, "Invalid time format, use RFC3339 (e.g., 2025-08-15T14:00:00+07:00)", http.StatusBadRequest)
		return
	}

	// Pass the automatically retrieved staffUserID to the service.
	schedule, err := h.service.CreateSchedule(r.Context(), req.BookingID, staffUserID, apptTime, req.Notes)
	if err != nil {
		http.Error(w, "Failed to create schedule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}
