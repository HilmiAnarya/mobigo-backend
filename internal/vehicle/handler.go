package vehicle

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Handler holds the dependencies for the vehicle handlers.
type Handler struct {
	service Service
}

// NewHandler creates a new instance of the vehicle handler.
func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes sets up the routing for the vehicle feature.
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// All vehicle routes will be under the /api/vehicles prefix.
	r := router.PathPrefix("/api/vehicles").Subrouter()

	r.HandleFunc("", h.createVehicleHandler).Methods("POST")
	r.HandleFunc("", h.getAllVehiclesHandler).Methods("GET")
	r.HandleFunc("/{id}", h.getVehicleByIDHandler).Methods("GET")
	r.HandleFunc("/{id}", h.updateVehicleHandler).Methods("PUT")
	r.HandleFunc("/{id}", h.deleteVehicleHandler).Methods("DELETE")
}

// createVehicleRequest defines the expected JSON body for creating a vehicle.
type createVehicleRequest struct {
	Make        string  `json:"make"`
	Model       string  `json:"model"`
	Year        int     `json:"year"`
	VIN         string  `json:"vin"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
}

// createVehicleHandler handles the creation of a new vehicle.
func (h *Handler) createVehicleHandler(w http.ResponseWriter, r *http.Request) {
	var req createVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.VIN == "" || req.Make == "" || req.Model == "" {
		http.Error(w, "VIN, Make, and Model are required", http.StatusBadRequest)
		return
	}

	vehicle, err := h.service.CreateVehicle(r.Context(), req.Make, req.Model, req.VIN, req.Description, req.Status, req.Year, req.Price)
	if err != nil {
		http.Error(w, "Failed to create vehicle", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicle)
}

// getAllVehiclesHandler handles retrieving all vehicles.
func (h *Handler) getAllVehiclesHandler(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.service.GetAllVehicles(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve vehicles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicles)
}

// getVehicleByIDHandler handles retrieving a single vehicle by its ID.
func (h *Handler) getVehicleByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	vehicle, err := h.service.GetVehicleByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to retrieve vehicle", http.StatusInternalServerError)
		return
	}
	if vehicle == nil {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicle)
}

// updateVehicleRequest defines the expected JSON body for updating a vehicle.
type updateVehicleRequest struct {
	Make        string  `json:"make"`
	Model       string  `json:"model"`
	Year        int     `json:"year"`
	VIN         string  `json:"vin"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
}

// updateVehicleHandler handles updating an existing vehicle.
func (h *Handler) updateVehicleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	var req updateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedVehicle, err := h.service.UpdateVehicle(r.Context(), id, req.Make, req.Model, req.VIN, req.Description, req.Status, req.Year, req.Price)
	if err != nil {
		http.Error(w, "Failed to update vehicle", http.StatusInternalServerError)
		return
	}
	if updatedVehicle == nil {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedVehicle)
}

// deleteVehicleHandler handles deleting a vehicle.
func (h *Handler) deleteVehicleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteVehicle(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete vehicle", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content is standard for successful deletes
}
