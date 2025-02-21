package testutil

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/jeevanions/italian-learning/internal/db/migrations"
	_ "modernc.org/sqlite"
)

// NewTestDB creates a temporary file-based SQLite database for testing
func NewTestDB() (*sql.DB, error) {
	// Create a temporary directory for the test database
	tmpDir, err := os.MkdirTemp("", "test-db-*")
	if err != nil {
		return nil, err
	}

	// Create database file path
	dbPath := filepath.Join(tmpDir, "test.db")

	// Open the database
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		return nil, err
	}

	// Apply migrations
	if err := migrations.Execute(db); err != nil {
		db.Close()
		os.RemoveAll(tmpDir)
		return nil, err
	}

	// Return database connection
	return db, nil
}
