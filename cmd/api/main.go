package main

import (
	"github.com/rs/cors"
	"log"
	"mobigo-backend/internal/user"
	"mobigo-backend/internal/vehicle"
	"net/http"
	"time"

	"mobigo-backend/pkg/database"
)

// apiHandlers is a container struct that holds all the different
// feature handlers for our application.
type apiHandlers struct {
	userHandler    *user.Handler
	vehicleHandler *vehicle.Handler
}

var jwtSecret = "a_very_secret_key_that_should_be_long_and_random"

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
	userRepository := user.NewGORMRepository(db)
	userService := user.NewService(userRepository, jwtSecret, 5*time.Second)
	userHandler := user.NewHandler(userService)

	vehicleRepository := vehicle.NewGORMRepository(db)
	vehicleService := vehicle.NewService(vehicleRepository, 5*time.Second)
	vehicleHandler := vehicle.NewHandler(vehicleService)

	// 3. Create the master handler container
	handlers := &apiHandlers{
		userHandler:    userHandler,
		vehicleHandler: vehicleHandler,
	}

	// 4. Define Routes
	router := defineRoutes(handlers)

	// 5. Setup CORS using rs/cors
	// This creates a new CORS handler with our desired options.
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Your React app's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// 6. Start Server with CORS middleware
	addr := ":8080"
	log.Printf("Server starting on %s\n", addr)
	// We wrap our main router with the CORS handler.
	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(addr, handler))
}
