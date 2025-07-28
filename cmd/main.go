package main

import (
	"aiload/internal"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	state, err := internal.NewState()
	if err != nil {
		log.Fatalf("failed to initialize state: %v", err)
	}
	defer state.DB.Close()
	defer state.Cache.Close()

	app := fiber.New()

	internal.RegisterRoutes(app, state)

	log.Fatal(app.Listen(":14001"))
}
