package main

import (
	"log"

	"aiload/internal/server"
	"aiload/internal/state"
)

func main() {
	// Initialize database
	db := state.InitDB("aiload.db")
	// Initialize cache
	cache := state.InitCache("aiload.db.cache")

	// Create a new Fiber application
	app := server.New(db, cache)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
