package repository

import (
	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

func (r *SQLiteRepository) GetAllStudySessions(limit, offset int) ([]models.StudySession, error) {
	query := `
		SELECT 
			id, 
			group_id,
			study_activity_id,
			created_at
		FROM study_sessions
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		err := rows.Scan(
			&session.ID,
			&session.StudyActivityID, // activity_id maps to StudyActivityID
			&session.GroupID,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *SQLiteRepository) GetTotalStudySessions() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&count)
	return count, err
}
