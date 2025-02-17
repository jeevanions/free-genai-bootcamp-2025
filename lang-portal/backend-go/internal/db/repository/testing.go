package repository

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func SetupTestDB(t *testing.T) *sql.DB {
	// Use in-memory SQLite for testing
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	// Create tables
	err = runTestMigrations(db)
	require.NoError(t, err)

	return db
}

func runTestMigrations(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE words (
			id INTEGER PRIMARY KEY,
			italian TEXT NOT NULL,
			english TEXT NOT NULL,
			parts_of_speech TEXT NOT NULL,
			gender TEXT CHECK (gender IN ('masculine', 'feminine', 'neuter') OR gender IS NULL),
			number TEXT CHECK (number IN ('singular', 'plural') OR number IS NULL),
			difficulty_level INTEGER NOT NULL CHECK (difficulty_level BETWEEN 1 AND 5),
			verb_conjugation JSON,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE groups (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description TEXT NOT NULL,
			difficulty_level INTEGER NOT NULL CHECK (difficulty_level BETWEEN 1 AND 5),
			category TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE study_activities (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL CHECK (type IN ('vocabulary', 'grammar', 'pronunciation')),
			requires_audio BOOLEAN NOT NULL DEFAULT FALSE,
			difficulty_level INTEGER NOT NULL CHECK (difficulty_level BETWEEN 1 AND 5),
			instructions TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE words_groups (
			word_id INTEGER NOT NULL,
			group_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (word_id, group_id),
			FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
		);

		CREATE TABLE study_sessions (
			id INTEGER PRIMARY KEY,
			group_id INTEGER NOT NULL,
			study_activity_id INTEGER NOT NULL,
			total_words INTEGER NOT NULL CHECK (total_words > 0),
			correct_words INTEGER NOT NULL CHECK (correct_words >= 0),
			duration_seconds INTEGER NOT NULL CHECK (duration_seconds > 0),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE RESTRICT,
			FOREIGN KEY (study_activity_id) REFERENCES study_activities(id) ON DELETE RESTRICT,
			CHECK (correct_words <= total_words)
		);

		CREATE TABLE word_review_items (
			id INTEGER PRIMARY KEY,
			word_id INTEGER NOT NULL,
			study_session_id INTEGER NOT NULL,
			correct BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
			FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE
		);
	`)
	return err
}
