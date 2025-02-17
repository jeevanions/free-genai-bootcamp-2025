//go:build mage

package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/jeevanions/italian-learning/internal/config"
	"github.com/jeevanions/italian-learning/internal/db/repository"
	"github.com/jeevanions/italian-learning/internal/db/seed"
	"github.com/jeevanions/italian-learning/internal/domain/services"
	"github.com/jeevanions/italian-learning/internal/pkg/backup"
	"github.com/magefile/mage/sh"
)

const (
	dbPath = "./data/italian-learning.db"
)

// ensureDataDir creates the data directory if it doesn't exist
func ensureDataDir() error {
	dir := filepath.Dir(dbPath)
	return os.MkdirAll(dir, 0755)
}

// cleanDB removes the existing database file if it exists
func cleanDB() error {
	if _, err := os.Stat(dbPath); err == nil {
		return os.Remove(dbPath)
	}
	return nil
}

// Generate runs code generation tools
func Generate() error {
	fmt.Println("Running code generation...")

	// Run sqlc
	if err := sh.Run("sqlc", "generate"); err != nil {
		return fmt.Errorf("sqlc generate failed: %w", err)
	}

	// Run swag
	if err := sh.Run("swag", "init", "-g", "cmd/server/main.go"); err != nil {
		return fmt.Errorf("swagger generation failed: %w", err)
	}

	return nil
}

// Reset removes and recreates the database
func Reset() error {
	fmt.Println("Resetting database...")
	if err := cleanDB(); err != nil {
		return fmt.Errorf("failed to clean database: %w", err)
	}
	return Migrate()
}

// Build builds the application
func Build() error {
	fmt.Println("Building application...")
	return sh.Run("go", "build", "-o", "bin/server", "./cmd/server")
}

// Dev runs the application in development mode
func Dev() error {
	fmt.Println("Running in development mode...")
	return sh.Run("go", "run", "./cmd/server")
}

// Test runs the test suite
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "./...", "-v", "-cover")
}

// Docs generates Swagger documentation
func Docs() error {
	fmt.Println("Generating Swagger documentation...")

	// Clean any existing docs
	if err := os.RemoveAll("./docs"); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to clean docs directory: %w", err)
		}
	}

	// Create docs directory
	if err := os.MkdirAll("./docs", 0755); err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	// Generate swagger docs
	if err := sh.Run("swag", "init",
		"-g", "cmd/server/main.go",
		"--parseDependency",
		"--parseInternal",
		"--parseDepth", "2",
		"--output", "docs"); err != nil {
		return fmt.Errorf("swagger generation failed: %w", err)
	}

	return nil
}

// Tidy updates Go module dependencies
func Tidy() error {
	fmt.Println("Updating dependencies...")
	return sh.Run("go", "mod", "tidy")
}

// Migrate runs database migrations
func Migrate() error {
	fmt.Println("Running database migrations...")
	if err := ensureDataDir(); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	return sh.Run("goose", "-dir", "internal/db/migrations", "sqlite3", dbPath, "up")
}

// MigrateDown rolls back database migrations
func MigrateDown() error {
	fmt.Println("Rolling back database migrations...")
	return sh.Run("goose", "-dir", "internal/db/migrations", "sqlite3", dbPath, "down")
}

// Seed adds initial data to the database
func Seed() error {
	fmt.Println("Seeding database...")

	// Initialize services
	db, err := setupDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	wordRepo := repository.NewSQLiteWordRepository(db)
	groupRepo := repository.NewSQLiteGroupRepository(db)

	wordService := services.NewWordService(wordRepo)
	groupService := services.NewGroupService(groupRepo)

	return seed.SeedBasicData(wordService, groupService)
}

func setupDatabase() (*sql.DB, error) {
	cfg := &config.Config{
		DatabasePath: "./data/italian-learning.db",
	}
	return config.NewDB(cfg)
}

// Backup creates a database backup
func Backup() error {
	fmt.Println("Creating database backup...")
	cfg := backup.NewBackupConfig()
	return backup.CreateBackup(cfg)
}

// BackupSchedule starts a scheduled backup process
func BackupSchedule() error {
	fmt.Println("Starting scheduled backups (Ctrl+C to stop)...")
	cfg := backup.NewBackupConfig()

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// Create initial backup
	if err := backup.CreateBackup(cfg); err != nil {
		return err
	}

	// Wait for ticker or interrupt
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			if err := backup.CreateBackup(cfg); err != nil {
				fmt.Printf("Backup failed: %v\n", err)
			}
		case <-sigChan:
			fmt.Println("\nStopping scheduled backups")
			return nil
		}
	}
}

// Clean removes generated files and artifacts
func Clean() error {
	fmt.Println("Cleaning generated files...")
	if err := cleanDB(); err != nil {
		return err
	}
	return os.RemoveAll("bin")
}
