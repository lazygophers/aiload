package state

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"go.etcd.io/bbolt"
)

// InitCache initializes the bbolt cache.
func InitCache(path string) *bbolt.DB {
	cache, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("failed to open cache: %v", err)
	}

	// Create a default bucket if it doesn't exist
	err = cache.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		return err
	})

	if err != nil {
		log.Fatalf("failed to create default bucket: %v", err)
	}

	log.Println("Cache database opened.")
	return cache
}

// NewTestCache creates a new temporary bbolt database for testing.
func NewTestCache(t testing.TB) (*bbolt.DB, func()) {
	t.Helper()

	// Create a temporary file for the bbolt database
	f, err := ioutil.TempFile(t.TempDir(), "test.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	path := f.Name()
	f.Close()

	cache, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		t.Fatalf("failed to open cache: %v", err)
	}

	// Create a default bucket if it doesn't exist
	err = cache.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		return err
	})

	if err != nil {
		t.Fatalf("failed to create default bucket: %v", err)
	}

	cleanup := func() {
		cache.Close()
	}

	return cache, cleanup
}
