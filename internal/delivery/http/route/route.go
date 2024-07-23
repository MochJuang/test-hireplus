package route

import (
	"github.com/gofiber/fiber/v2"
	"hireplus-project/internal/config"
	httpdelivery "hireplus-project/internal/delivery/http"
	middleware "hireplus-project/internal/delivery/http/midlleware"
	"hireplus-project/internal/service"
)

func SetupRoutes(app *fiber.App, userService service.UserService, transactionService service.TransactionService, cfg config.Config) {
	// Initialize http
	userController := httpdelivery.NewUserController(userService)
	transactionController := httpdelivery.NewTransactionController(transactionService)
	app.Use(middleware.ErrorHandlerMiddleware)

	// Public routes
	app.Post("/api/register", userController.Register)
	app.Post("/api/login", userController.Login)

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware(cfg))
	api.Post("/topup", transactionController.TopUp)
	api.Post("/pay", transactionController.Payment)
	api.Post("/transfer", transactionController.Transfer)
	api.Get("/transactions", transactionController.TransactionsReport)
	api.Put("/profile", userController.UpdateProfile)
}
