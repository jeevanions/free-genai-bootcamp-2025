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
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return activity, nil
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
