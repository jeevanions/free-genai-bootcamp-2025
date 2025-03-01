package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"
)

//go:embed *.sql
var migrationFiles embed.FS

func Execute(db *sql.DB) error {
	// Debug: Print available files
	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// Debug: Print found files
	log.Printf("Found migration files:")
	for _, entry := range entries {
		log.Printf("- %s", entry.Name())
	}

	// Read all SQL files from the embedded filesystem
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	// Execute each migration file
	for _, file := range files {
		content, err := migrationFiles.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Split the file into Up and Down migrations
		parts := strings.Split(string(content), "-- +goose Down")
		if len(parts) != 2 {
			return fmt.Errorf("invalid migration format in %s", file)
		}

		// Extract the Up migration (remove the -- +goose Up line)
		upMigration := strings.TrimPrefix(parts[0], "-- +goose Up")
		upMigration = strings.TrimSpace(upMigration)

		// Execute the Up migration
		_, err = db.Exec(upMigration)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	return nil
}
