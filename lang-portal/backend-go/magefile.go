//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Build builds the application
func Build() error {
	fmt.Println("Building...")
	// Update dependencies
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to update dependencies: %w", err)
	}
	// Generate Swagger docs
	if err := sh.Run("swag", "init", "-g", "cmd/server/main.go"); err != nil {
		return fmt.Errorf("failed to generate swagger docs: %w", err)
	}
	return sh.Run("go", "build", "-o", "bin/server", "./cmd/server")
}

// Run runs the application
func Run() error {
	fmt.Println("Running...")
	return sh.Run("go", "run", "./cmd/server/main.go")
}

// Test runs the tests
func Test() error {
	// Update dependencies
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to update dependencies: %w", err)
	}
	// Unit tests
	sh.Run("go", "test", "-short", "./...")

	// Integration tests
	sh.Run("go", "test", "./test/integration/...")

	// E2E tests
	sh.Run("k6", "run", "./test/e2e/...")

	return nil
}

// Docs generates Swagger documentation
func Docs() error {
	fmt.Println("Generating Swagger documentation...")
	// Update dependencies
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to update dependencies: %w", err)
	}
	return sh.Run("swag", "init", "-g", "cmd/server/main.go")
}

// Tidy updates Go module dependencies
func Tidy() error {
	fmt.Println("Updating dependencies...")
	return sh.Run("go", "mod", "tidy")
}
