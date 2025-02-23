package repository

import (
	"database/sql"
	"time"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

func (r *SQLiteRepository) GetStudyActivities(limit, offset int) (*models.StudyActivityListResponse, error) {
	query := `
		SELECT id, name, thumbnail_url, description, created_at
		FROM study_activities
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.StudyActivityResponse
	for rows.Next() {
		var activity models.StudyActivityResponse
		var thumbnailURL, description sql.NullString
		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&thumbnailURL,
			&description,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if thumbnailURL.Valid {
			url := thumbnailURL.String
			activity.ThumbnailURL = &url
		}
		if description.Valid {
			desc := description.String
			activity.Description = &desc
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Get total count for pagination
	var total int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM study_activities").Scan(&total); err != nil {
		return nil, err
	}

	return &models.StudyActivityListResponse{
		Items: activities,
		Pagination: models.PaginationResponse{
			CurrentPage:  offset/limit + 1,
			TotalPages:   (total + limit - 1) / limit,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	}, nil
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
