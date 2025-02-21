//go:build mage
// +build mage

package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/jeevanions/lang-portal/backend-go/internal/db/seeder"
	_ "modernc.org/sqlite"
)

const (
	dbFile = "words.db"
)

// Default target to run when none is specified
var Default = Build

// Build builds the application
func Build() error {
	fmt.Println("Building...")
	if err := os.MkdirAll("tmp", 0755); err != nil {
		return err
	}
	return sh.Run("go", "build", "-o", "tmp/main", "./cmd/server")
}

// Clean cleans up build artifacts
func Clean() error {
	fmt.Println("Cleaning...")
	dirs := []string{"bin", "tmp", "docs"}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

// Test runs the test suite
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "-v", "./...")
}

// Coverage runs tests with coverage
func Coverage() error {
	fmt.Println("Running tests with coverage...")
	if err := sh.Run("go", "test", "-coverprofile=coverage.out", "./..."); err != nil {
		return err
	}
	return sh.Run("go", "tool", "cover", "-html=coverage.out")
}

// Lint runs golangci-lint
func Lint() error {
	fmt.Println("Running linter...")
	return sh.Run("golangci-lint", "run")
}

// Generate runs code generation tools
func Generate() error {
	mg.SerialDeps(SQLc, Swagger)
	fmt.Println("Running go generate...")
	return sh.Run("go", "generate", "./...")
}

// SQLc generates database code
func SQLc() error {
	fmt.Println("Generating SQLc code...")
	cmd := exec.Command("/Users/jeevan/go/bin/sqlc", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Swagger generates API documentation
func Swagger() error {
	fmt.Println("Generating Swagger documentation...")
	cmd := exec.Command("/Users/jeevan/go/bin/swag", "init", "-g", "cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type DB mg.Namespace

// DB:Migrate runs database migrations
func (DB) Migrate() error {
	fmt.Println("Running database migrations...")
	return sh.Run("goose", "-dir", "internal/db/migrations", "sqlite3", dbFile, "up")
}

// DB:Seed seeds the database with initial data
func (DB) Seed() error {
	fmt.Println("Seeding database...")
	
	// Ensure migrations are up to date
	mg.SerialDeps(DB.Migrate)

	// Open database connection
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Create seeder instance
	seeder := seeder.New(db)

	// Run seeder
	if err := seeder.SeedFromJSON("internal/db/seeds"); err != nil {
		return fmt.Errorf("failed to seed database: %w", err)
	}

	fmt.Println("Database seeded successfully")
	return nil
}

// DB:Reset resets the database
func (DB) Reset() error {
	fmt.Println("Resetting database...")
	if err := os.Remove(dbFile); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// DB:Status shows migration status
func (DB) Status() error {
	fmt.Println("Checking migration status...")
	return sh.Run("goose", "-dir", "internal/db/migrations", "sqlite3", dbFile, "status")
}

// Run runs the application
func Run() error {
	fmt.Println("Running server...")
	cmd := exec.Command("go", "run", "./cmd/server")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Dev runs the application in development mode with live reload
func Dev() error {
	fmt.Println("Running server in development mode...")
	cmd := exec.Command("/Users/jeevan/go/bin/air")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Install installs development dependencies
func Install() error {
	fmt.Println("Installing development dependencies...")
	deps := []string{
		"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
		"github.com/pressly/goose/v3/cmd/goose@v3.19.0",
		"github.com/swaggo/swag/cmd/swag@v1.16.3",
		"github.com/air-verse/air@latest",
		"github.com/sqlc-dev/sqlc/cmd/sqlc@v1.24.0",
	}

	for _, dep := range deps {
		fmt.Printf("Installing %s...\n", dep)
		cmd := exec.Command("go", "install", dep)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install %s: %w", dep, err)
		}
	}
	return nil
}

// All runs all tasks in sequence: clean, generate, build, test
func All() {
	mg.SerialDeps(Clean, Generate, Build, Test)
}
