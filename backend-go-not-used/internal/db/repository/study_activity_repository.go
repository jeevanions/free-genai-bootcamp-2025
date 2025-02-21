package repository

import (
	"context"
	"database/sql"

	"github.com/jeevanions/italian-learning/internal/db/models"
)

type SQLiteStudyActivityRepository struct {
	db *sql.DB
}

func NewSQLiteStudyActivityRepository(db *sql.DB) *SQLiteStudyActivityRepository {
	return &SQLiteStudyActivityRepository{db: db}
}

func (r *SQLiteStudyActivityRepository) Create(ctx context.Context, activity *models.StudyActivity) error {
	query := `
		INSERT INTO study_activities (
			name, type, requires_audio, difficulty_level, instructions
		) VALUES (?, ?, ?, ?, ?)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		activity.Name,
		activity.Type,
		activity.RequiresAudio,
		activity.DifficultyLevel,
		activity.Instructions,
	).Scan(&activity.ID, &activity.CreatedAt)

	return err
}

func (r *SQLiteStudyActivityRepository) GetByID(ctx context.Context, id int64) (*models.StudyActivity, error) {
	activity := &models.StudyActivity{}
	query := `
		SELECT id, name, type, requires_audio, difficulty_level, 
			   instructions, created_at
		FROM study_activities
		WHERE id = ?`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&activity.ID,
		&activity.Name,
		&activity.Type,
		&activity.RequiresAudio,
		&activity.DifficultyLevel,
		&activity.Instructions,
		&activity.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return activity, err
}

func (r *SQLiteStudyActivityRepository) GetCategories(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT type
		FROM study_activities
		ORDER BY type`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *SQLiteStudyActivityRepository) GetRecommended(ctx context.Context, count int) ([]*models.StudyActivity, error) {
	// Get activities that match the user's current level and haven't been completed yet
	query := `
		SELECT id, name, type, requires_audio, difficulty_level, instructions, created_at
		FROM study_activities a
		WHERE NOT EXISTS (
			SELECT 1 FROM study_sessions s
			WHERE s.study_activity_id = a.id
		)
		ORDER BY difficulty_level ASC
		LIMIT ?`

	rows, err := r.db.QueryContext(ctx, query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*models.StudyActivity
	for rows.Next() {
		activity := &models.StudyActivity{}
		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&activity.Type,
			&activity.RequiresAudio,
			&activity.DifficultyLevel,
			&activity.Instructions,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *SQLiteStudyActivityRepository) List(ctx context.Context, offset, limit int) ([]*models.StudyActivity, error) {
	query := `
		SELECT id, name, type, requires_audio, difficulty_level, 
			   instructions, created_at
		FROM study_activities
		ORDER BY difficulty_level, name
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []*models.StudyActivity
	for rows.Next() {
		activity := &models.StudyActivity{}
		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&activity.Type,
			&activity.RequiresAudio,
			&activity.DifficultyLevel,
			&activity.Instructions,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, rows.Err()
}
