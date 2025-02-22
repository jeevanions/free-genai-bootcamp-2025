package repository

func (r *SQLiteRepository) ResetHistory() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete all study sessions and word reviews
	_, err = tx.Exec("DELETE FROM word_reviews")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		return err
	}

	// Reset word statistics
	_, err = tx.Exec("UPDATE words SET correct_count = 0, wrong_count = 0")
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
		"word_reviews",
		"study_sessions",
		"words",
		"groups",
		"study_activities",
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
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    italian TEXT NOT NULL,
    english TEXT NOT NULL,
    parts JSON,
    correct_count INTEGER DEFAULT 0,
    wrong_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS study_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS study_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    activity_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (activity_id) REFERENCES study_activities(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE IF NOT EXISTS word_reviews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    word_id INTEGER NOT NULL,
    is_correct BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES study_sessions(id),
    FOREIGN KEY (word_id) REFERENCES words(id)
);`, nil
}
