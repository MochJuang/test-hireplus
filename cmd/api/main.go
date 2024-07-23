package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"hireplus-project/internal/config"
	"hireplus-project/internal/delivery/http/route"
	"hireplus-project/internal/repository/mysql"
	"hireplus-project/internal/service"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Initialize database
	db, err := mysql.NewConnector(cfg)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Perform database migration
	err = mysql.Migrate(db)
	if err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}

	// Initialize repositories
	userRepo := mysql.NewUserRepository(db)
	transactionRepo := mysql.NewTransactionRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, cfg)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)

	// Initialize Fiber app
	app := fiber.New()

	app.Use(logger.New())

	// Setup routes
	route.SetupRoutes(app, userService, transactionService, cfg)

	// Setup error handler middleware

	// Start server
	log.Fatal(app.Listen(cfg.ServerAddress))
}
