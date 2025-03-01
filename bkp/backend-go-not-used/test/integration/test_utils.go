package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeevanions/italian-learning/internal/api/handlers"
	"github.com/jeevanions/italian-learning/internal/api/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
	"github.com/jeevanions/italian-learning/internal/domain/services"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Create a temporary database file
	tmpfile, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)

	// Clean up the file when we're done
	t.Cleanup(func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	})

	// Open the database
	db, err := sql.Open("sqlite", tmpfile.Name())
	require.NoError(t, err)

	// Create tables
	_, err = db.ExecContext(context.Background(), `
		CREATE TABLE words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			italian TEXT NOT NULL,
			english TEXT NOT NULL,
			parts_of_speech TEXT NOT NULL,
			gender TEXT,
			number TEXT,
			difficulty_level INTEGER NOT NULL,
			verb_conjugation JSON,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(italian),
			CHECK (gender IS NULL OR gender IN ('masculine', 'feminine', 'neuter')),
			CHECK (number IS NULL OR number IN ('singular', 'plural')),
			CHECK (difficulty_level BETWEEN 1 AND 5),
			CHECK (parts_of_speech IN ('noun', 'verb', 'adjective', 'adverb', 'preposition', 'conjunction', 'interjection'))
		);

		CREATE INDEX idx_words_italian ON words(italian);
		CREATE INDEX idx_words_difficulty ON words(difficulty_level);
		CREATE INDEX idx_words_parts_speech ON words(parts_of_speech);
	`)
	require.NoError(t, err)

	return db
}

func setupTestServer(t *testing.T, db *sql.DB) *gin.Engine {
	// Create repository
	wordRepo := repository.NewSQLiteWordRepository(db)

	// Create service
	wordService := services.NewWordService(wordRepo)

	// Create handler
	wordHandler := handlers.NewWordHandler(wordService)

	// Setup router
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		words := v1.Group("/words")
		{
			words.GET("/", wordHandler.ListWords)
			words.POST("/", wordHandler.CreateWord)
			words.GET("/search", wordHandler.SearchWords)
			words.GET("/filters", wordHandler.GetFilters)
			words.GET("/:id", wordHandler.GetWord)
		}
	}

	return router
}

func createTestActivity(t *testing.T, router *gin.Engine) *models.StudyActivity {
	activity := models.CreateStudyActivityRequest{
		Name:            "Test Activity",
		Type:            "vocabulary",
		DifficultyLevel: 1,
		RequiresAudio:   false,
		Instructions:    "Test instructions",
	}

	body, err := json.Marshal(activity)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/activities/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var response models.StudyActivity
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.NotZero(t, response.ID)

	return &response
}
