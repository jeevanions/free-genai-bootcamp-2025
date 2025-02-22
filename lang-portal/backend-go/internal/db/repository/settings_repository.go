package repository

func (r *SQLiteRepository) ResetHistory() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete all study sessions and word reviews
	_, err = tx.Exec("DELETE FROM word_review_items")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		return err
	}

	// Try to reset word statistics if columns exist
	_, err = tx.Exec(`
		UPDATE words 
		SET correct_count = CASE 
			WHEN (SELECT COUNT(*) FROM pragma_table_info('words') WHERE name='correct_count') > 0 THEN 0 
			ELSE correct_count 
			END,
			wrong_count = CASE 
			WHEN (SELECT COUNT(*) FROM pragma_table_info('words') WHERE name='wrong_count') > 0 THEN 0 
			ELSE wrong_count 
			END
	`)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SQLiteRepository) DropAllTables() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// List of tables to drop
	tables := []string{
		"word_review_items",
		"study_sessions",
		"words_groups",
		"study_activities",
		"groups",
		"words",
	}

	// Drop each table if it exists
	for _, table := range tables {
		_, err := tx.Exec(`DROP TABLE IF EXISTS ` + table)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLiteRepository) CreateTables() error {
	// Read the schema file
	schema, err := r.readSchemaFile()
	if err != nil {
		return err
	}

	// Execute the schema
	_, err = r.db.Exec(schema)
	return err
}

func (r *SQLiteRepository) readSchemaFile() (string, error) {
	return `
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    italian TEXT NOT NULL,
    english TEXT NOT NULL,
    parts JSON,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    words_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS words_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS study_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    thumbnail_url TEXT,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS study_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    study_activity_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    FOREIGN KEY (study_activity_id) REFERENCES study_activities(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS word_review_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word_id INTEGER NOT NULL,
    study_session_id INTEGER NOT NULL,
    correct BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_words_groups_word_id ON words_groups(word_id);
CREATE INDEX IF NOT EXISTS idx_words_groups_group_id ON words_groups(group_id);
CREATE INDEX IF NOT EXISTS idx_study_sessions_group_id ON study_sessions(group_id);
CREATE INDEX IF NOT EXISTS idx_word_review_items_word_id ON word_review_items(word_id);
CREATE INDEX IF NOT EXISTS idx_word_review_items_study_session_id ON word_review_items(study_session_id);
`, nil
}
