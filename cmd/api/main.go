package main

import (
	"log"
	"mobigo-backend/internal/user"
	"net/http"
	"time"

	"mobigo-backend/pkg/database"
)

// apiHandlers is a container struct that holds all the different
// feature handlers for our application.
type apiHandlers struct {
	userHandler *user.Handler
	// vehicleHandler *vehicle.Handler // We will add this later
}

func main() {
	// 1. Initialize Database Connection
	dbUser := "root"
	dbPassword := "" // Use your MySQL root password
	dbName := "mobigo_db"

	db, err := database.Connect(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	log.Println("Successfully connected to the database using GORM!")

	// 2. Wire up all dependencies
	// Create instances of every repository, service, and handler.
	userRepository := user.NewGORMRepository(db)
	userService := user.NewService(userRepository, 5*time.Second)
	userHandler := user.NewHandler(userService)

	// vehicleRepository := vehicle.NewGORMRepository(db)
	// vehicleService := vehicle.NewService(vehicleRepository, 5*time.Second)
	// vehicleHandler := vehicle.NewHandler(vehicleService)

	// 3. Create the master handler container
	// This single object holds all our handlers.
	handlers := &apiHandlers{
		userHandler: userHandler,
		// vehicleHandler: vehicleHandler,
	}

	// 4. Define Routes, passing the handler container
	router := defineRoutes(handlers)

	// Placeholder health check route
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is healthy and running on Clean Architecture with GORM"))
	}).Methods("GET")

	// 5. Start Server
	addr := ":8080"
	log.Printf("Server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
