package repository

import (
	"context"
	"database/sql"

	"github.com/jeevanions/italian-learning/internal/db/models"
)

type SQLiteGroupRepository struct {
	db *sql.DB
}

func NewSQLiteGroupRepository(db *sql.DB) *SQLiteGroupRepository {
	return &SQLiteGroupRepository{db: db}
}

func (r *SQLiteGroupRepository) Create(ctx context.Context, group *models.Group) error {
	query := `
		INSERT INTO groups (
			name, description, difficulty_level, category
		) VALUES (?, ?, ?, ?)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		group.Name,
		group.Description,
		group.DifficultyLevel,
		group.Category,
	).Scan(&group.ID, &group.CreatedAt)

	return err
}

func (r *SQLiteGroupRepository) GetByID(ctx context.Context, id int64) (*models.Group, error) {
	group := &models.Group{}
	query := `
		SELECT id, name, description, difficulty_level, category, created_at
		FROM groups
		WHERE id = ?`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&group.ID,
		&group.Name,
		&group.Description,
		&group.DifficultyLevel,
		&group.Category,
		&group.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (r *SQLiteGroupRepository) GetWordCount(ctx context.Context, groupID int64) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM words_groups
		WHERE group_id = ?`

	err := r.db.QueryRowContext(ctx, query, groupID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SQLiteGroupRepository) AddWordsToGroup(ctx context.Context, groupID int64, wordIDs []int64) (int, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Prepare insert statement
	query := `
		INSERT INTO words_groups (word_id, group_id, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT (word_id, group_id) DO NOTHING`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	// Insert each word and count successful inserts
	wordsAdded := 0
	for _, wordID := range wordIDs {
		result, err := stmt.ExecContext(ctx, wordID, groupID)
		if err != nil {
			return 0, err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		wordsAdded += int(rowsAffected)
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return wordsAdded, nil
}

func (r *SQLiteGroupRepository) GetGroupStatistics(ctx context.Context, groupID int64) (apimodels.GroupStatistics, error) {
	stats := apimodels.GroupStatistics{}
	query := `
		WITH StudyStats AS (
			SELECT 
				COUNT(DISTINCT ws.word_id) as studied_words,
				AVG(CASE WHEN ws.is_correct THEN 1.0 ELSE 0.0 END) as success_rate,
				COUNT(DISTINCT ss.id) as total_sessions,
				AVG(ss.duration_seconds) as avg_duration
			FROM study_sessions ss
			JOIN word_study ws ON ws.session_id = ss.id
			JOIN words_groups wg ON wg.word_id = ws.word_id
			WHERE wg.group_id = ?
		)
		SELECT 
			(SELECT COUNT(*) FROM words_groups WHERE group_id = ?) as total_words,
			COALESCE(studied_words, 0),
			COALESCE(success_rate, 0.0),
			COALESCE(total_sessions, 0),
			COALESCE(avg_duration, 0)
		FROM StudyStats`

	err := r.db.QueryRowContext(ctx, query, groupID, groupID).Scan(
		&stats.TotalWords,
		&stats.StudiedWords,
		&stats.SuccessRate,
		&stats.TotalSessions,
		&stats.AverageDuration,
	)
	if err != nil && err != sql.ErrNoRows {
		return stats, err
	}

	return stats, nil
}

func (r *SQLiteGroupRepository) GetGroupProgress(ctx context.Context, groupID int64) (apimodels.GroupProgress, error) {
	progress := apimodels.GroupProgress{}
	query := `
		WITH ProgressStats AS (
			SELECT 
				AVG(CASE WHEN ws.mastery_level >= 3 THEN 1.0 ELSE 0.0 END) as mastery_percentage,
				MAX(ss.created_at) as last_study_date,
				(SELECT COUNT(DISTINCT DATE(created_at))
				FROM study_sessions
				WHERE id IN (
					SELECT DISTINCT ss2.id
					FROM study_sessions ss2
					JOIN word_study ws2 ON ws2.session_id = ss2.id
					JOIN words_groups wg2 ON wg2.word_id = ws2.word_id
					WHERE wg2.group_id = ?
					AND ss2.created_at >= DATE('now', '-30 days')
				)
			) as study_streak
			FROM study_sessions ss
			JOIN word_study ws ON ws.session_id = ss.id
			JOIN words_groups wg ON wg.word_id = ws.word_id
			WHERE wg.group_id = ?
		)
		SELECT 
			COALESCE(mastery_percentage, 0.0),
			COALESCE(last_study_date, ''),
			COALESCE(study_streak, 0)
		FROM ProgressStats`

	err := r.db.QueryRowContext(ctx, query, groupID, groupID).Scan(
		&progress.MasteryPercentage,
		&progress.LastStudyDate,
		&progress.Streak,
	)
	if err != nil && err != sql.ErrNoRows {
		return progress, err
	}

	return progress, nil
}

func (r *SQLiteGroupRepository) List(ctx context.Context, offset, limit int) ([]*models.Group, error) {
	query := `
		SELECT id, name, description, difficulty_level, category, created_at
		FROM groups
		ORDER BY difficulty_level, name
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		group := &models.Group{}
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.Description,
			&group.DifficultyLevel,
			&group.Category,
			&group.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, rows.Err()
}

func (r *SQLiteGroupRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM groups WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *SQLiteGroupRepository) Update(ctx context.Context, group *models.Group) error {
	query := `UPDATE groups 
			  SET name = ?, description = ?, difficulty_level = ?, category = ? 
			  WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query,
		group.Name, group.Description, group.DifficultyLevel, group.Category, group.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
