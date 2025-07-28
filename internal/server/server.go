package server

import (
	"aiload/internal/handler"
	"aiload/internal/repository"
	"aiload/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.etcd.io/bbolt"
	"gorm.io/gorm"
)

// New creates a new Fiber application.
func New(db *gorm.DB, cache *bbolt.DB) *fiber.App {
	app := fiber.New()

	// Use logger middleware
	app.Use(logger.New())

	// Create dependencies
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cache)
	userHandler := handler.NewUserHandler(userService)

	// Setup routes
	SetupRoutes(app, userHandler)

	return app
}

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	// API group
	api := app.Group("/api")

	// Ping endpoint
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "pong",
		})
	})

	// User routes
	user := api.Group("/users")
	user.Post("/", userHandler.CreateUser)
	user.Get("/:id", userHandler.GetUser)
	user.Put("/:id", userHandler.UpdateUser)
	user.Delete("/:id", userHandler.DeleteUser)
}
