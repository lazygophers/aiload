package internal

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"go.etcd.io/bbolt"
)

// State holds the application's state, including database connections.
type State struct {
	DB    *sql.DB
	Cache *bbolt.DB
}

// NewState creates and initializes a new State object.
func NewState() (*State, error) {
	// Ensure the data directory exists.
	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Connect to the SQLite database.
	db, err := sql.Open("sqlite", "./data/aiload.db")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sqlite: %w", err)
	}

	// Open the bbolt cache database.
	cache, err := bbolt.Open("./data/aiload.bbolt", 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open bbolt cache: %w", err)
	}

	return &State{
		DB:    db,
		Cache: cache,
	}, nil
}
