package seeder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
)

type Seeder struct {
	db *repository.SQLiteRepository
}

func New(db *repository.SQLiteRepository) *Seeder {
	return &Seeder{db: db}
}

type Word struct {
	ID        int64           `json:"id"`
	Italian   string         `json:"italian"`
	English   string         `json:"english"`
	Parts     json.RawMessage `json:"parts"`
	CreatedAt string         `json:"created_at"`
}

type Group struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	WordsCount int    `json:"words_count"`
	Level      string `json:"level"`
}

type WordGroup struct {
	WordID  int64 `json:"word_id"`
	GroupID int64 `json:"group_id"`
}

type StudyActivity struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
}

func (s *Seeder) SeedFromJSON(seedDir string) error {
	// Begin transaction
	tx, err := s.db.DB().Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Seed groups
	groups, err := s.loadGroups(filepath.Join(seedDir, "groups.json"))
	if err != nil {
		return fmt.Errorf("failed to load groups: %w", err)
	}

	for _, group := range groups {
		_, err := tx.Exec(
			"INSERT INTO groups (name) VALUES (?)",
			group.Name,
		)
		if err != nil {
			return fmt.Errorf("failed to insert group %s: %w", group.Name, err)
		}
	}

	// Seed words
	words, err := s.loadWords(filepath.Join(seedDir, "words.json"))
	if err != nil {
		return fmt.Errorf("failed to load words: %w", err)
	}

	for _, word := range words {
		_, err := tx.Exec(
			"INSERT INTO words (italian, english, parts) VALUES (?, ?, ?)",
			word.Italian, word.English, word.Parts,
		)
		if err != nil {
			return fmt.Errorf("failed to insert word %s: %w", word.Italian, err)
		}
	}

	// Seed word_groups
	wordGroups, err := s.loadWordGroups(filepath.Join(seedDir, "words_groups.json"))
	if err != nil {
		return fmt.Errorf("failed to load word groups: %w", err)
	}

	for _, wg := range wordGroups {
		_, err := tx.Exec(
			"INSERT INTO words_groups (word_id, group_id) VALUES (?, ?)",
			wg.WordID, wg.GroupID,
		)
		if err != nil {
			return fmt.Errorf("failed to insert word group mapping: %w", err)
		}
	}

	// Seed study activities
	activities, err := s.loadStudyActivities(filepath.Join(seedDir, "study_activities.json"))
	if err != nil {
		return fmt.Errorf("failed to load study activities: %w", err)
	}

	for _, activity := range activities {
		_, err := tx.Exec(
			"INSERT INTO study_activities (name, thumbnail_url, description) VALUES (?, ?, ?)",
			activity.Name, activity.ThumbnailURL, activity.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to insert study activity %s: %w", activity.Name, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Seeder) loadGroups(path string) ([]Group, error) {
	var data struct {
		Groups []Group `json:"groups"`
	}
	if err := loadJSON(path, &data); err != nil {
		return nil, err
	}
	return data.Groups, nil
}

func (s *Seeder) loadWords(path string) ([]Word, error) {
	var data struct {
		Words []Word `json:"words"`
	}
	if err := loadJSON(path, &data); err != nil {
		return nil, err
	}
	return data.Words, nil
}

func (s *Seeder) loadWordGroups(path string) ([]WordGroup, error) {
	var data struct {
		WordGroups []WordGroup `json:"words_groups"`
	}
	if err := loadJSON(path, &data); err != nil {
		return nil, err
	}
	return data.WordGroups, nil
}

func (s *Seeder) loadStudyActivities(path string) ([]StudyActivity, error) {
	var data struct {
		Activities []StudyActivity `json:"study_activities"`
	}
	if err := loadJSON(path, &data); err != nil {
		return nil, err
	}
	return data.Activities, nil
}

func loadJSON(path string, v interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from %s: %w", path, err)
	}

	return nil
}
