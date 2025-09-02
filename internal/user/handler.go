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
	staffRouter.HandleFunc("/login", h.loginStaffHandler).Methods("POST") // Add the new login route

	// Customer routes
	customerRouter := router.PathPrefix("/api/customers").Subrouter()
	customerRouter.HandleFunc("/register", h.registerCustomerHandler).Methods("POST")
	// We can reuse the login handler for customers, as the logic is identical.
	customerRouter.HandleFunc("/login", h.loginCustomerHandler).Methods("POST")
}

// --- General Structs (used by both staff and customer) ---

// FIX: Renamed to be general-purpose.
type registrationResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// --- Customer Handlers ---

type registerCustomerRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

func (h *Handler) registerCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var req registerCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		http.Error(w, "Full name, email, and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.service.RegisterCustomer(r.Context(), req.FullName, req.Email, req.Password, req.PhoneNumber)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Use the general response struct.
	resp := registrationResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) loginCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	token, err := h.service.LoginCustomer(r.Context(), req.Email, req.Password)
	if err != nil {
		// Handle specific errors
		if err.Error() == "access denied: staff cannot log in through customer portal" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err.Error() == "invalid email or password" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}

// --- Staff Handlers ---
type registerStaffRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
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
	resp := registrationResponse{
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

// loginStaffHandler handles the staff login request.
func (h *Handler) loginStaffHandler(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	token, err := h.service.LoginStaff(r.Context(), req.Email, req.Password)
	if err != nil {
		// Handle our new authorization error specifically
		if err.Error() == "access denied: user is not a staff member" {
			http.Error(w, err.Error(), http.StatusForbidden) // 403 Forbidden is the correct code
			return
		}
		if err.Error() == "invalid email or password" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}
