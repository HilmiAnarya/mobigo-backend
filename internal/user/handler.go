package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Handler holds the dependencies for the user handlers.
type Handler struct {
	service Service
}

// NewHandler creates a new instance of the user handler.
func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes sets up the routing for the user feature.
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// We are creating a subrouter for staff-related endpoints for better organization.
	staffRouter := router.PathPrefix("/api/staff").Subrouter()
	staffRouter.HandleFunc("/register", h.registerStaffHandler).Methods("POST")
}

// registerStaffRequest defines the expected JSON body for the registration request.
type registerStaffRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
}

// registerStaffResponse defines the JSON response for a successful registration.
type registerStaffResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
}

// registerStaffHandler handles the staff registration request.
func (h *Handler) registerStaffHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Decode the request body.
	var req registerStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Basic validation (can be expanded later).
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		http.Error(w, "Full name, email, and password are required", http.StatusBadRequest)
		return
	}

	// 3. Call the business logic (the service).
	// We pass the request context, which is important for managing request lifecycle.
	user, err := h.service.RegisterStaff(r.Context(), req.FullName, req.Email, req.Password, req.PhoneNumber, req.Address)
	if err != nil {
		// Check for specific business errors we defined in the service.
		if err.Error() == "user with this email already exists" {
			http.Error(w, err.Error(), http.StatusConflict) // 409 Conflict
			return
		}
		// For any other error, it's a server-side problem.
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// 4. Create the response.
	resp := registerStaffResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
	}

	// 5. Send the successful JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(resp)
}
