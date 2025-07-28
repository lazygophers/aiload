package internal

import "github.com/gofiber/fiber/v2"

// RegisterRoutes sets up the routes for the application.
func RegisterRoutes(app *fiber.App, state *State) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
