package main

import (
	"log"
	"net/http"

	"mobigo-backend/pkg/database"

	"github.com/gorilla/mux"
)

func main() {
	// 1. Initialize Database Connection
	// IMPORTANT: Replace with your actual MySQL username and password.
	dbUser := "root"
	dbPassword := "" // Use your MySQL root password if you have one
	dbName := "mobigo_db"

	db, err := database.Connect(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to the database!")

	// 2. Setup Router
	router := mux.NewRouter()

	// 3. Wire up dependencies (we will do this in the next checkpoints)
	// e.g., userRepository := user.NewMySQLRepository(db)
	// e.g., userUsecase := user.NewUsecase(userRepository)
	// e.g., user.RegisterHandlers(router, userUsecase)

	// Placeholder health check route
	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API is healthy and running on Clean Architecture"))
	}).Methods("GET")

	// 4. Start Server
	addr := ":8080"
	log.Printf("Server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
