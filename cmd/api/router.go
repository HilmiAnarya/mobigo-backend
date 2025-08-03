package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// defineRoutes now accepts the apiHandlers container, giving it access to all handlers.
func defineRoutes(handlers *apiHandlers) *mux.Router {
	router := mux.NewRouter()

	// Register routes for the user feature using the userHandler from the container.
	handlers.userHandler.RegisterRoutes(router)

	// When we add the vehicle feature, we will just add one line here:
	handlers.vehicleHandler.RegisterRoutes(router)

	handlers.bookingHandler.RegisterRoutes(router) // Register the new booking routes

	// General-purpose routes
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is healthy and running on Clean Architecture with GORM"))
	}).Methods("GET")

	return router
}
