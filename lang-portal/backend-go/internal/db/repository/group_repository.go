package repository

import (
	"database/sql"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

func (r *SQLiteRepository) GetGroups(limit, offset int) (*models.GroupListResponse, error) {
	query := `
		SELECT 
			g.id, g.name,
			(SELECT COUNT(*) FROM words_groups wg WHERE wg.group_id = g.id) as word_count
		FROM groups g
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.GroupResponse
	for rows.Next() {
		var group models.GroupResponse
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.WordCount,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM groups"
	var total int
	err = r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &models.GroupListResponse{
		Items: groups,
		Pagination: models.PaginationResponse{
			CurrentPage:  offset/limit + 1,
			TotalPages:   (total + limit - 1) / limit,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	}, nil
}

func (r *SQLiteRepository) GetGroupByID(id int64) (*models.GroupDetailResponse, error) {
	query := `
		SELECT 
			g.id, g.name,
			(SELECT COUNT(*) FROM words_groups wg WHERE wg.group_id = g.id) as word_count
		FROM groups g
		WHERE g.id = ?
	`

	var group models.GroupDetailResponse
	err := r.db.QueryRow(query, id).Scan(
		&group.ID,
		&group.Name,
		&group.Stats.TotalWordCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *SQLiteRepository) GetGroupWords(groupID int64, limit, offset int) (*models.GroupWordsResponse, error) {
	query := `
		SELECT 
			w.id, w.italian, w.english, w.parts,
			COUNT(CASE WHEN wri.correct THEN 1 END) as correct_count,
			COUNT(CASE WHEN NOT wri.correct THEN 1 END) as wrong_count
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wg.group_id = ?
		GROUP BY w.id
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []models.WordResponse
	for rows.Next() {
		var word models.WordResponse
		var partsStr string
		err := rows.Scan(
			&word.ID,
			&word.Italian,
			&word.English,
			&partsStr,
			&word.CorrectCount,
			&word.WrongCount,
		)
		if err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		WHERE wg.group_id = ?
	`
	var total int
	err = r.db.QueryRow(countQuery, groupID).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &models.GroupWordsResponse{
		Items: words,
		Pagination: models.PaginationResponse{
			CurrentPage:  offset/limit + 1,
			TotalPages:   (total + limit - 1) / limit,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	}, nil
}

func (r *SQLiteRepository) GetGroupStudySessions(groupID int64, limit, offset int) (*models.GroupStudySessionsResponse, error) {
	query := `
		SELECT 
			ss.id,
			sa.name as activity_name,
			g.name as group_name,
			ss.created_at,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.study_session_id = ss.id) as words_count,
			(SELECT COUNT(*) FROM word_review_items wri WHERE wri.study_session_id = ss.id AND wri.correct) as correct_count
		FROM study_sessions ss
		JOIN study_activities sa ON ss.study_activity_id = sa.id
		JOIN groups g ON ss.group_id = g.id
		WHERE ss.group_id = ?
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.StudySessionResponse
	for rows.Next() {
		var session models.StudySessionResponse
		err := rows.Scan(
			&session.ID,
			&session.ActivityName,
			&session.GroupName,
			&session.CreatedAt,
			&session.WordsCount,
			&session.CorrectCount,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM study_sessions WHERE group_id = ?"
	var total int
	err = r.db.QueryRow(countQuery, groupID).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &models.GroupStudySessionsResponse{
		Items: sessions,
		Pagination: models.PaginationResponse{
			CurrentPage:  offset/limit + 1,
			TotalPages:   (total + limit - 1) / limit,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	}, nil
}
