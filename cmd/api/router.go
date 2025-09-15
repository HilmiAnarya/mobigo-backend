package main

import (
	"mobigo-backend/pkg/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// defineRoutes now accepts the apiHandlers container, giving it access to all handlers.
func defineRoutes(handlers *apiHandlers, jwtSecret string) *mux.Router {
	router := mux.NewRouter()

	authMiddleware := middleware.JWTAuthMiddleware(jwtSecret)

	// --- Serve Static Files ---
	// This is crucial. It creates a route that allows the frontend to access
	// the images saved in the ./uploads directory.
	fs := http.FileServer(http.Dir("./uploads/"))
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", fs))

	// Pass the middleware to the handlers that need it
	handlers.userHandler.RegisterRoutes(router)
	handlers.vehicleHandler.RegisterRoutes(router, authMiddleware)
	handlers.bookingHandler.RegisterRoutes(router, authMiddleware)
	handlers.scheduleHandler.RegisterRoutes(router, authMiddleware)
	handlers.agreementHandler.RegisterRoutes(router, authMiddleware)
	handlers.paymentHandler.RegisterRoutes(router, authMiddleware)
	handlers.vehicleImageHandler.RegisterRoutes(router, authMiddleware) // This registers all image-related routes

	// General-purpose routes
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is healthy and running on Clean Architecture with GORM"))
	}).Methods("GET")

	return router
}
