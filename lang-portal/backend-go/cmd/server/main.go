package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jeevanions/lang-portal/backend-go/internal/db/seeder"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/jeevanions/lang-portal/backend-go/internal/api/router"
	"github.com/jeevanions/lang-portal/backend-go/internal/config"
	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
	_ "github.com/jeevanions/lang-portal/backend-go/docs" // Import swagger docs
)

// @title Italian Language Learning Portal API
// @version 1.0
// @description API for the Italian Language Learning Portal
// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	// Initialize logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := repository.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	defer db.Close()

	// Initialize seeder
	seeder := seeder.New(db)

	// Initialize router
	r := router.Setup(db, seeder)

	// Initialize HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}
