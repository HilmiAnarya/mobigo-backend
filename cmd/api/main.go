package main

import (
	"github.com/robfig/cron/v3"
	"github.com/rs/cors"
	"log"
	"mobigo-backend/internal/agreement"
	"mobigo-backend/internal/booking"
	"mobigo-backend/internal/installment"
	"mobigo-backend/internal/payment"
	"mobigo-backend/internal/schedule"
	"mobigo-backend/internal/task"
	"mobigo-backend/internal/user"
	"mobigo-backend/internal/vehicle"
	"mobigo-backend/internal/vehicleimage"
	"net/http"
	"time"

	"mobigo-backend/pkg/database"
)

// apiHandlers is a container struct that holds all the different
// feature handlers for our application.
type apiHandlers struct {
	userHandler         *user.Handler
	vehicleHandler      *vehicle.Handler
	bookingHandler      *booking.Handler
	scheduleHandler     *schedule.Handler
	agreementHandler    *agreement.Handler
	paymentHandler      *payment.Handler
	vehicleImageHandler *vehicleimage.Handler // Add the vehicle image handler
}

func main() {
	// 1. Initialize Database Connection
	dbUser := "root"
	dbPassword := "" // Use your MySQL root password
	dbName := "mobigo_db"
	var jwtSecret = "a_very_secret_key_that_should_be_long_and_random"

	db, err := database.Connect(dbUser, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	log.Println("Successfully connected to the database using GORM!")

	// --- Dependency Injection (Wiring) ---
	// Build repositories
	userRepository := user.NewGORMRepository(db)
	vehicleRepository := vehicle.NewGORMRepository(db)
	bookingRepository := booking.NewGORMRepository(db)
	scheduleRepository := schedule.NewGORMRepository(db)
	agreementRepository := agreement.NewGORMRepository(db)
	paymentRepository := payment.NewGORMRepository(db)
	installmentRepository := installment.NewGORMRepository(db)
	vehicleImageRepository := vehicleimage.NewGORMRepository(db) // New repository

	// Build services
	userService := user.NewService(userRepository, jwtSecret, 5*time.Second)
	vehicleService := vehicle.NewService(vehicleRepository, 5*time.Second)
	bookingService := booking.NewService(bookingRepository, scheduleRepository, vehicleRepository)
	scheduleService := schedule.NewService(scheduleRepository)
	agreementService := agreement.NewService(agreementRepository, bookingRepository)
	paymentService := payment.NewService(paymentRepository, installmentRepository, vehicleRepository, agreementRepository, bookingRepository)
	vehicleImageService := vehicleimage.NewService(vehicleImageRepository) // New service

	// Build handlers
	userHandler := user.NewHandler(userService)
	vehicleHandler := vehicle.NewHandler(vehicleService)
	bookingHandler := booking.NewHandler(bookingService)
	scheduleHandler := schedule.NewHandler(scheduleService)
	agreementHandler := agreement.NewHandler(agreementService, paymentService)
	paymentHandler := payment.NewHandler(paymentService)
	vehicleImageHandler := vehicleimage.NewHandler(vehicleImageService) // New handler

	// 3. Create the master handler container
	handlers := &apiHandlers{
		userHandler:         userHandler,
		vehicleHandler:      vehicleHandler,
		bookingHandler:      bookingHandler,
		scheduleHandler:     scheduleHandler,
		agreementHandler:    agreementHandler,
		paymentHandler:      paymentHandler,
		vehicleImageHandler: vehicleImageHandler, // Add handler to the container
	}

	// 4. Define Routes
	router := defineRoutes(handlers, jwtSecret)

	// --- Setup Cron Jobs ---
	c := cron.New()
	penaltyChecker := task.NewPenaltyChecker(installmentRepository)
	// THE FIX: Set the schedule to run once a day at midnight ("0 0 * * *").
	_, err = c.AddFunc("1 * * * *", penaltyChecker.Run)
	if err != nil {
		log.Fatalf("Could not add cron job: %v", err)
	}
	c.Start()
	log.Println("Cron job scheduler started. Penalty check will run daily at midnight.")
	defer c.Stop()

	// 5. Setup CORS using rs/cors
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Your React app's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// 6. Start Server with CORS middleware
	addr := ":8080"
	log.Printf("Server starting on %s\n", addr)
	// We wrap our main router with the CORS handler.
	handler := corsHandler.Handler(router)
	log.Fatal(http.ListenAndServe(addr, handler))
}
