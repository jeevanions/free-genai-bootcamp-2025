package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

func (r *SQLiteRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *SQLiteRepository) GetGroupIDByName(name string) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT id FROM groups WHERE name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SQLiteRepository) CreateGroup(name string) (int64, error) {
	result, err := r.db.Exec("INSERT INTO groups (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *SQLiteRepository) CreateWord(word *models.WordResponse) (int64, error) {
	partsJSON, err := json.Marshal(word.Parts)
	if err != nil {
		return 0, err
	}

	result, err := r.db.Exec(
		"INSERT INTO words (italian, english, parts) VALUES (?, ?, ?)",
		word.Italian,
		word.English,
		partsJSON,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *SQLiteRepository) AddWordToGroup(wordID, groupID int64) error {
	_, err := r.db.Exec(
		"INSERT INTO words_groups (word_id, group_id) VALUES (?, ?)",
		wordID,
		groupID,
	)
	return err
}

func (r *SQLiteRepository) UpdateGroupWordsCount(groupID int64) error {
	_, err := r.db.Exec(`
		UPDATE groups 
		SET words_count = (
			SELECT COUNT(*) 
			FROM words_groups 
			WHERE group_id = ?
		)
		WHERE id = ?`,
		groupID,
		groupID,
	)
	return err
}
