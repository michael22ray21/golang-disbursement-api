// Paper.id Wallet Disbursement API
// This API provides a single endpoint for disbursing funds from a user's wallet to their designated bank account.

package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

// Config struct for application configuration
type Config struct {
	Port           string
	DSN            string
	AllowedOrigins []string
}

// App struct for application dependencies
type App struct {
	Config    Config
	DB        *sql.DB
	Validator *validator.Validate
}

func main() {
	// Initialize config
	config := Config{
		Port:           ":8080",
		DSN:            "./paper_id.db",
		AllowedOrigins: []string{"*"},
	}

	// get env port if available
	if port := os.Getenv("PORT"); port != "" {
		config.Port = ":" + port
	}

	// Initialize database
	db, err := initializeDB(config.DSN)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// Initialize application
	app := &App{
		Config:    config,
		DB:        db,
		Validator: validator.New(),
	}

	// Seed test data
	if err := seedTestData(app); err != nil {
		log.Printf("Warning: Error seeding test data: %v", err)
	}

	//^ start of router
	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))

	// Register routes
	r.Get("/", app.handleIndex)
	r.Post("/api/disbursements", app.handleDisbursement)
	//^ end of router

	// CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Start server
	log.Printf("Starting server on port %s", config.Port)
	err = http.ListenAndServe(config.Port, corsMiddleware.Handler(r))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
