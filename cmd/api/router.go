package main

import (
	"mobigo-backend/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// defineRoutes now accepts the apiHandlers container, giving it access to all handlers.
func defineRoutes(handlers *apiHandlers) *mux.Router {
	router := mux.NewRouter()

	authMiddleware := middleware.JWTAuthMiddleware(jwtSecret)

	// Pass the middleware to the handlers that need it
	handlers.userHandler.RegisterRoutes(router)                    // User routes don't need protection
	handlers.vehicleHandler.RegisterRoutes(router, authMiddleware) // Protect vehicle routes
	handlers.bookingHandler.RegisterRoutes(router, authMiddleware) // Protect booking routes

	// General-purpose routes
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is healthy and running on Clean Architecture with GORM"))
	}).Methods("GET")

	return router
}
