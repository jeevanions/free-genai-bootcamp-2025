package repository

import (
	"time"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

func (r *SQLiteRepository) GetStudyActivities(limit, offset int) ([]models.StudyActivity, error) {
	query := `
		SELECT id, name, thumbnail_url, description, launch_url, created_at
		FROM study_activities
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.StudyActivity
	for rows.Next() {
		var activity models.StudyActivity
		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&activity.ThumbnailURL,
			&activity.Description,
			&activity.LaunchURL,
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

func (r *SQLiteRepository) CreateStudyActivitySession(activityID, groupID int64) (*models.LaunchStudyActivityResponse, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert study session
	sessionQuery := `
		INSERT INTO study_sessions (group_id, created_at)
		VALUES (?, datetime('now'))
	`
	sessionResult, err := tx.Exec(sessionQuery, groupID)
	if err != nil {
		return nil, err
	}

	sessionID, err := sessionResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Insert study activity
	activityQuery := `
		INSERT INTO study_activities (study_session_id, group_id, created_at)
		VALUES (?, ?, datetime('now'))
	`
	activityResult, err := tx.Exec(activityQuery, sessionID, groupID)
	if err != nil {
		return nil, err
	}

	activityID, err = activityResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &models.LaunchStudyActivityResponse{
		StudySessionID:  sessionID,
		StudyActivityID: activityID,
		GroupID:        groupID,
		CreatedAt:      time.Now(),
	}, nil
}
