package state

import (
	"aiload/internal/model"
	"log"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes the database connection.
func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate schema
	db.AutoMigrate(&model.User{})
	log.Println("Database connection established.")
	return db
}

// NewTestDB creates a new in-memory database for testing.
func NewTestDB(t testing.TB) (*gorm.DB, func()) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate schema
	db.AutoMigrate(&model.User{})

	cleanup := func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("failed to get sqlDB: %v", err)
		}
		sqlDB.Close()
	}

	return db, cleanup
}
