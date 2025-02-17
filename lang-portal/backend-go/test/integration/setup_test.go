package integration

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/db/migrations"
	"github.com/jeevanions/italian-learning/internal/api/router"
	"github.com/jeevanions/italian-learning/internal/db/repository"
	"github.com/jeevanions/italian-learning/internal/domain/services"
	"github.com/jeevanions/italian-learning/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestApplicationStartup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Use test server instead of real server
	r, _, err := setupTestServer()
	require.NoError(t, err)
	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("Health Check", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/api/v1/health")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestGracefulShutdown(t *testing.T) {
	// Setup test server
	gin.SetMode(gin.TestMode)
	logger.Setup()

	r, _, err := setupTestServer()
	require.NoError(t, err)
	srv := &http.Server{
		Addr:    ":8082",
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Give the server time to start
	time.Sleep(100 * time.Millisecond)

	// Make a test request
	_, err = http.Get("http://localhost:8082/health")
	require.NoError(t, err)

	// Initiate shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	assert.NoError(t, err, "Server shutdown failed")

	// Verify server has shut down
	_, err = http.Get("http://localhost:8082/health")
	assert.Error(t, err, "Expected error after shutdown")
}

func setupTestServer() (*gin.Engine, *sql.DB, error) {
	// Initialize in-memory SQLite for testing
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, nil, err
	}

	// Run migrations using the actual migration system
	if err := migrations.Execute(db); err != nil {
		return nil, nil, err
	}

	// Initialize repositories
	wordRepo := repository.NewSQLiteWordRepository(db)
	groupRepo := repository.NewSQLiteGroupRepository(db)
	sessionRepo := repository.NewSQLiteStudySessionRepository(db)
	reviewRepo := repository.NewSQLiteWordReviewRepository(db)
	activityRepo := repository.NewSQLiteStudyActivityRepository(db)

	// Initialize services
	wordService := services.NewWordService(wordRepo)
	groupService := services.NewGroupService(groupRepo)
	sessionService := services.NewStudySessionService(sessionRepo, groupRepo, activityRepo)
	reviewService := services.NewWordReviewService(reviewRepo, wordRepo, sessionRepo)
	activityService := services.NewStudyActivityService(activityRepo)

	// Setup Gin router
	r := gin.Default()
	v1 := r.Group("/api/v1")
	router.SetupRoutes(v1, wordService, groupService, sessionService, reviewService, activityService)

	return r, db, nil
}

func setupTestDB(t *testing.T) *sql.DB {
	// Initialize in-memory SQLite for testing
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	// Run migrations using the actual migration system
	err = migrations.Execute(db)
	require.NoError(t, err)

	return db
}

func runMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY,
			italian TEXT NOT NULL,
			english TEXT NOT NULL,
			parts_of_speech TEXT NOT NULL,
			gender TEXT,
			number TEXT,
			difficulty_level INTEGER NOT NULL,
			verb_conjugation JSON,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT NOT NULL,
			difficulty_level INTEGER NOT NULL,
			category TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS study_activities (
			id INTEGER PRIMARY KEY,
			type TEXT NOT NULL,
			difficulty_level INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS study_sessions (
			id INTEGER PRIMARY KEY,
			group_id INTEGER REFERENCES groups(id),
			study_activity_id INTEGER REFERENCES study_activities(id),
			total_words INTEGER NOT NULL,
			correct_words INTEGER NOT NULL,
			duration_seconds INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS word_review_items (
			id INTEGER PRIMARY KEY,
			word_id INTEGER REFERENCES words(id),
			study_session_id INTEGER REFERENCES study_sessions(id),
			correct BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS words_groups (
			word_id INTEGER REFERENCES words(id),
			group_id INTEGER REFERENCES groups(id),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (word_id, group_id)
		);
	`)
	return err
}
