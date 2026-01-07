package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"movie-booking/api/v1"
	"movie-booking/api/v1/controllers"
	"movie-booking/api/v1/middleware"
	"movie-booking/config"
	"movie-booking/core/services"
	coretypes "movie-booking/core/types"
	"movie-booking/datastore"
	"movie-booking/dbmigrations"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Parse flags
	apiFlag := flag.Bool("api", false, "Start API server")
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	migrationCommand := flag.String("migration-command", "up", "Migration command: up, down, or status")
	flag.Parse()

	// Initialize config
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// Initialize logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations if requested
	if *migrateFlag {
		if err := dbmigrations.RunMigrations(db, *migrationCommand); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations completed successfully")
		return
	}

	// Start API server if requested
	if *apiFlag {
		startAPIServer(db)
		return
	}

	// Default: show usage
	fmt.Println("Usage:")
	fmt.Println("  --api                    Start API server")
	fmt.Println("  --migrate                Run database migrations")
	fmt.Println("  --migration-command=up    Migration command (up, down, status)")
}

func initDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetDatabaseUser(),
		config.GetDatabasePassword(),
		config.GetDatabaseHost(),
		config.GetDatabasePort(),
		config.GetDatabaseName(),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.Info("Database connection established")
	return db, nil
}

func startAPIServer(db *gorm.DB) {
	// Create datastore
	store := datastore.NewDataStore(db)

	// Create clients (for dependency injection)
	clients := &coretypes.Clients{}

	// Create services
	authService := services.NewAuthService(clients, store)
	movieService := services.NewMovieService(clients, store)
	showService := services.NewShowService(clients, store)
	seatService := services.NewSeatService(clients, store)
	bookingService := services.NewBookingService(clients, store)

	// Create controller
	ctrl := controllers.NewController(
		authService,
		movieService,
		showService,
		seatService,
		bookingService,
	)

	// Create router
	router := mux.NewRouter()

	// Add health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	// Register API routes
	if err := v1.AddRoutesToRouter(router, ctrl); err != nil {
		log.Fatalf("Failed to register routes: %v", err)
	}

	// Add CORS middleware - wrap the router
	handler := middleware.CORS(router)
	
	// Also add a catch-all OPTIONS handler for CORS preflight
	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This will be handled by CORS middleware, but ensure it doesn't 405
	})

	// Start server
	port := config.GetServerPort()
	logrus.WithField("port", port).Info("Starting API server")
	
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
