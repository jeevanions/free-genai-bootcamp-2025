package repository

import (
	"context"
	"database/sql"

	"github.com/jeevanions/italian-learning/internal/db/models"
)

type SQLiteWordRepository struct {
	db *sql.DB
}

func NewSQLiteWordRepository(db *sql.DB) *SQLiteWordRepository {
	return &SQLiteWordRepository{db: db}
}

func (r *SQLiteWordRepository) Create(ctx context.Context, word *models.Word) error {
	query := `
		INSERT INTO words (
			italian, english, parts_of_speech, gender, number,
			difficulty_level, verb_conjugation, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		word.Italian,
		word.English,
		word.PartsOfSpeech,
		sql.NullString{String: word.Gender.String, Valid: word.Gender.Valid},
		sql.NullString{String: word.Number.String, Valid: word.Number.Valid},
		word.DifficultyLevel,
		word.VerbConjugation,
		sql.NullString{String: word.Notes.String, Valid: word.Notes.Valid},
	).Scan(&word.ID, &word.CreatedAt)

	return err
}

func (r *SQLiteWordRepository) GetByID(ctx context.Context, id int64) (*models.Word, error) {
	word := &models.Word{}
	var verbConjugation sql.NullString

	query := `SELECT id, italian, english, parts_of_speech, gender, number,
			  difficulty_level, verb_conjugation, notes, created_at 
			  FROM words WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&word.ID,
		&word.Italian,
		&word.English,
		&word.PartsOfSpeech,
		&word.Gender,
		&word.Number,
		&word.DifficultyLevel,
		&verbConjugation,
		&word.Notes,
		&word.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Only set VerbConjugation if the database value is not NULL
	if verbConjugation.Valid {
		word.VerbConjugation = []byte(verbConjugation.String)
	}

	return word, nil
}

func (r *SQLiteWordRepository) List(ctx context.Context, offset, limit int) ([]*models.Word, error) {
	query := `
		SELECT id, italian, english, parts_of_speech, gender, number,
			   difficulty_level, verb_conjugation, notes, created_at
		FROM words
		ORDER BY id
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*models.Word
	for rows.Next() {
		word := &models.Word{}
		var verbConjugation sql.NullString
		err := rows.Scan(
			&word.ID,
			&word.Italian,
			&word.English,
			&word.PartsOfSpeech,
			&word.Gender,
			&word.Number,
			&word.DifficultyLevel,
			&verbConjugation,
			&word.Notes,
			&word.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if verbConjugation.Valid {
			word.VerbConjugation = []byte(verbConjugation.String)
		}
		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func (r *SQLiteWordRepository) Search(ctx context.Context, query string, offset, limit int) ([]*models.Word, error) {
	sqlQuery := `
		SELECT id, italian, english, parts_of_speech, gender, number,
			   difficulty_level, verb_conjugation, notes, created_at
		FROM words
		WHERE italian LIKE ? COLLATE NOCASE 
		   OR english LIKE ? COLLATE NOCASE 
		   OR notes LIKE ? COLLATE NOCASE
		ORDER BY 
			CASE 
				WHEN italian LIKE ? THEN 1
				WHEN english LIKE ? THEN 2
				ELSE 3
			END,
			id
		LIMIT ? OFFSET ?`

	likeQuery := "%" + query + "%"
	rows, err := r.db.QueryContext(ctx, sqlQuery, 
		likeQuery, likeQuery, likeQuery,
		likeQuery, likeQuery,
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []*models.Word
	for rows.Next() {
		word := &models.Word{}
		var verbConjugation sql.NullString
		err := rows.Scan(
			&word.ID,
			&word.Italian,
			&word.English,
			&word.PartsOfSpeech,
			&word.Gender,
			&word.Number,
			&word.DifficultyLevel,
			&verbConjugation,
			&word.Notes,
			&word.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if verbConjugation.Valid {
			word.VerbConjugation = []byte(verbConjugation.String)
		}
		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func (r *SQLiteWordRepository) SearchCount(ctx context.Context, query string) (int, error) {
	sqlQuery := `
		SELECT COUNT(*)
		FROM words
		WHERE italian LIKE ? COLLATE NOCASE 
		   OR english LIKE ? COLLATE NOCASE 
		   OR notes LIKE ? COLLATE NOCASE`

	likeQuery := "%" + query + "%"
	var count int
	err := r.db.QueryRowContext(ctx, sqlQuery, likeQuery, likeQuery, likeQuery).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SQLiteWordRepository) Count(ctx context.Context) (int, error) {
	sqlQuery := `SELECT COUNT(*) FROM words`

	var count int
	err := r.db.QueryRowContext(ctx, sqlQuery).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SQLiteWordRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM words WHERE id = ?`
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

func (r *SQLiteWordRepository) Update(ctx context.Context, word *models.Word) error {
	query := `UPDATE words SET 
		italian = ?, english = ?, parts_of_speech = ?, 
		gender = ?, number = ?, difficulty_level = ?,
		verb_conjugation = ?, notes = ?
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		word.Italian, word.English, word.PartsOfSpeech,
		sql.NullString{String: word.Gender.String, Valid: word.Gender.Valid},
		sql.NullString{String: word.Number.String, Valid: word.Number.Valid},
		word.DifficultyLevel,
		word.VerbConjugation,
		sql.NullString{String: word.Notes.String, Valid: word.Notes.Valid},
		word.ID,
	)
	return err
}
