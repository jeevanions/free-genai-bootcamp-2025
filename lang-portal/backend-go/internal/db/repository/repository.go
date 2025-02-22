package repository

import (
	"database/sql"
	"time"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
	_ "modernc.org/sqlite"
)

type Repository interface {
	// Dashboard queries
	GetLastStudySession() (*models.DashboardLastStudySession, error)
	GetStudyProgress() (*models.DashboardStudyProgress, error)
	GetQuickStats() (*models.DashboardQuickStats, error)

	// Study activities
	GetStudyActivities(limit, offset int) ([]models.StudyActivity, error)
	GetStudyActivity(id int64) (*models.StudyActivity, error)
	GetStudyActivitySessions(activityID int64, limit, offset int) ([]models.StudySession, error)
	CreateStudyActivitySession(activityID, groupID int64) (*models.LaunchStudyActivityResponse, error)
	GetWordReviewsBySessionID(sessionID int64) ([]models.WordReviewItem, error)

	// Words
	GetWords(limit, offset int) (*models.WordListResponse, error)
	GetWordByID(id int64) (*models.WordResponse, error)

	// Groups
	GetGroups(limit, offset int) (*models.GroupListResponse, error)
	GetGroupByID(id int64) (*models.GroupDetailResponse, error)
	GetGroupWords(groupID int64, limit, offset int) (*models.GroupWordsResponse, error)
	GetGroupStudySessions(groupID int64, limit, offset int) (*models.GroupStudySessionsResponse, error)

	// Study Sessions
	GetAllStudySessions(limit, offset int) ([]models.StudySession, error)
	GetTotalStudySessions() (int, error)

	// Close the database connection
	// Settings
	ResetHistory() error
	DropAllTables() error
	CreateTables() error

	// Close the database connection
	Close() error
}

type SQLiteRepository struct {
	db *sql.DB
}

func (r *SQLiteRepository) DB() *sql.DB {
	return r.db
}

func NewDB(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite only supports one writer at a time
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}

func (r *SQLiteRepository) GetLastStudySession() (*models.DashboardLastStudySession, error) {
	query := `
		SELECT 
			s.id,
			s.group_id,
			s.created_at,
			s.study_activity_id,
			g.name as group_name
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		ORDER BY s.created_at DESC
		LIMIT 1
	`

	var session models.DashboardLastStudySession
	err := r.db.QueryRow(query).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&session.StudyActivityID,
		&session.GroupName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SQLiteRepository) GetStudyProgress() (*models.DashboardStudyProgress, error) {
	query := `
		SELECT 
			(SELECT COUNT(DISTINCT word_id) FROM word_review_items) as total_words_studied,
			(SELECT COUNT(*) FROM words) as total_available_words
	`

	var progress models.DashboardStudyProgress
	err := r.db.QueryRow(query).Scan(
		&progress.TotalWordsStudied,
		&progress.TotalAvailableWords,
	)
	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func (r *SQLiteRepository) GetQuickStats() (*models.DashboardQuickStats, error) {
	query := `
		SELECT
			COALESCE(
				(
					SELECT CAST(SUM(CASE WHEN correct THEN 1 ELSE 0 END) AS FLOAT) * 100.0 / COUNT(*)
					FROM word_review_items
					WHERE created_at >= date('now', '-30 days')
				), 0.0
			) as success_rate,
			COALESCE((SELECT COUNT(DISTINCT id) FROM study_sessions), 0) as total_sessions,
			COALESCE(
				(SELECT COUNT(DISTINCT group_id)
				FROM study_sessions
				WHERE created_at >= date('now', '-30 days')
				), 0
			) as active_groups,
			COALESCE(
				(SELECT COUNT(DISTINCT date(created_at))
				FROM study_sessions
				WHERE created_at >= date('now', '-30 days')
				), 0
			) as streak_days
	`

	var stats models.DashboardQuickStats
	err := r.db.QueryRow(query).Scan(
		&stats.SuccessRate,
		&stats.TotalStudySessions,
		&stats.TotalActiveGroups,
		&stats.StudyStreakDays,
	)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *SQLiteRepository) GetStudyActivity(id int64) (*models.StudyActivity, error) {
	query := `
		SELECT id, name, thumbnail_url, description, created_at
		FROM study_activities
		WHERE id = ?
	`

	var activity models.StudyActivity
	err := r.db.QueryRow(query, id).Scan(
		&activity.ID,
		&activity.Name,
		&activity.ThumbnailURL,
		&activity.Description,
		&activity.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &activity, nil
}

func (r *SQLiteRepository) GetStudyActivitySessions(activityID int64, limit, offset int) ([]models.StudySession, error) {
	query := `
		SELECT id, group_id, study_activity_id, created_at
		FROM study_sessions
		WHERE study_activity_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, activityID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		err := rows.Scan(
			&session.ID,
			&session.GroupID,
			&session.StudyActivityID,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

func (r *SQLiteRepository) GetWordReviewsBySessionID(sessionID int64) ([]models.WordReviewItem, error) {
	query := `
		SELECT id, word_id, study_session_id, correct, created_at
		FROM word_review_items
		WHERE study_session_id = ?
	`

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.WordReviewItem
	for rows.Next() {
		var review models.WordReviewItem
		err := rows.Scan(
			&review.ID,
			&review.WordID,
			&review.StudySessionID,
			&review.Correct,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, rows.Err()
}


