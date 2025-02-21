package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/jeevanions/lang-portal/backend-go/internal/domain/models"
)

func (r *SQLiteRepository) GetWords(limit, offset int) (*models.WordListResponse, error) {
	query := `
		SELECT 
			w.id, w.italian, w.english, w.parts,
			COUNT(CASE WHEN wri.correct THEN 1 END) as correct_count,
			COUNT(CASE WHEN NOT wri.correct THEN 1 END) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		GROUP BY w.id
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
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
		// Parse parts JSON
		if err := json.Unmarshal([]byte(partsStr), &word.Parts); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM words"
	var total int
	err = r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &models.WordListResponse{
		Items: words,
		Pagination: models.PaginationResponse{
			CurrentPage:  offset/limit + 1,
			TotalPages:   (total + limit - 1) / limit,
			TotalItems:   total,
			ItemsPerPage: limit,
		},
	}, nil
}

func (r *SQLiteRepository) GetWordByID(id int64) (*models.WordResponse, error) {
	query := `
		SELECT 
			w.id, w.italian, w.english, w.parts,
			COUNT(CASE WHEN wri.correct THEN 1 END) as correct_count,
			COUNT(CASE WHEN NOT wri.correct THEN 1 END) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE w.id = ?
		GROUP BY w.id
	`

	var word models.WordResponse
	var partsStr string
	err := r.db.QueryRow(query, id).Scan(
		&word.ID,
		&word.Italian,
		&word.English,
		&partsStr,
		&word.CorrectCount,
		&word.WrongCount,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse parts JSON
	if err := json.Unmarshal([]byte(partsStr), &word.Parts); err != nil {
		return nil, err
	}

	return &word, nil
}
