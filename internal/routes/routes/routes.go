package routes

import (
	"github.com/notmkw/log/internal/handlers"
	"github.com/notmkw/log/internal/middleware"
	"github.com/notmkw/log/internal/repositories"
	"github.com/notmkw/log/internal/services"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	logger := middleware.NewEarlyRequestLogData(false)
	app.Use(middleware.LoggerMiddleware(logger))

	api := app.Group("/api/v1")

	api.Get("/health", userHandler.HealthCheck)

	userRoutes := api.Group("/user")
	userRoutes.Get("/:id", userHandler.GetUser)
	userRoutes.Post("/", userHandler.CreateUser)
	userRoutes.Put("/:id", userHandler.UpdateUser)
	userRoutes.Delete("/:id", userHandler.DeleteUser)
}
