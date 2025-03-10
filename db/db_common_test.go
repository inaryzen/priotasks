package db

import (
	"os"
	"testing"
)

func setupTestDB(t *testing.T) *DbSQLite {
	// Create a temporary database file
	tmpfile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}
	tmpfile.Close()

	// Initialize test database
	db := NewDbSQLite()
	db.Init(tmpfile.Name())

	// Return the database and cleanup function
	t.Cleanup(func() {
		db.Close()
		os.Remove(tmpfile.Name())
	})

	return db
}
