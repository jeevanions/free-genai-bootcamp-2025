package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/jeevanions/italian-learning/internal/api/models"
	dbmodels "github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInvalid  = errors.New("invalid input")
)

// WordService defines the interface for word-related operations
type WordService interface {
	CreateWord(ctx context.Context, req *models.CreateWordRequest) (*models.Word, error)
	GetWord(ctx context.Context, id int64) (*models.Word, error)
	ListWords(ctx context.Context, page, pageSize int) ([]models.Word, error)
}

// wordServiceImpl implements WordService
type wordServiceImpl struct {
	repo repository.WordRepository
}

// NewWordService creates a new WordService instance
func NewWordService(repo repository.WordRepository) WordService {
	return &wordServiceImpl{
		repo: repo,
	}
}

func (s *wordServiceImpl) CreateWord(ctx context.Context, req *models.CreateWordRequest) (*models.Word, error) {
	// Convert API request to DB model
	var dbVerbConj json.RawMessage
	if req.VerbConjugation != nil {
		dbVerbConj = json.RawMessage(*req.VerbConjugation)
	}

	dbWord := &dbmodels.Word{
		Italian:         req.Italian,
		English:         req.English,
		PartsOfSpeech:   req.PartsOfSpeech,
		DifficultyLevel: req.DifficultyLevel,
		VerbConjugation: dbVerbConj,
	}
	if req.Gender != nil {
		dbWord.Gender = sql.NullString{String: *req.Gender, Valid: true}
	}
	if req.Number != nil {
		dbWord.Number = sql.NullString{String: *req.Number, Valid: true}
	}
	if req.Notes != nil {
		dbWord.Notes = sql.NullString{String: *req.Notes, Valid: true}
	}

	// Create in DB
	if err := s.repo.Create(ctx, dbWord); err != nil {
		return nil, err
	}

	// Convert back to API model
	var apiVerbConj *string
	if len(dbWord.VerbConjugation) > 0 {
		s := string(dbWord.VerbConjugation)
		apiVerbConj = &s
	}

	return &models.Word{
		ID:              dbWord.ID,
		Italian:         dbWord.Italian,
		English:         dbWord.English,
		PartsOfSpeech:   dbWord.PartsOfSpeech,
		Gender:          req.Gender,
		Number:          req.Number,
		DifficultyLevel: dbWord.DifficultyLevel,
		VerbConjugation: apiVerbConj,
		Notes:           req.Notes,
		CreatedAt:       dbWord.CreatedAt,
	}, nil
}

func (s *wordServiceImpl) GetWord(ctx context.Context, id int64) (*models.Word, error) {
	dbWord, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Convert DB model to API model
	var gender, number, notes *string
	if dbWord.Gender.Valid {
		gender = &dbWord.Gender.String
	}
	if dbWord.Number.Valid {
		number = &dbWord.Number.String
	}
	if dbWord.Notes.Valid {
		notes = &dbWord.Notes.String
	}

	var verbConj *string
	if len(dbWord.VerbConjugation) > 0 {
		s := string(dbWord.VerbConjugation)
		verbConj = &s
	}

	return &models.Word{
		ID:              dbWord.ID,
		Italian:         dbWord.Italian,
		English:         dbWord.English,
		PartsOfSpeech:   dbWord.PartsOfSpeech,
		Gender:          gender,
		Number:          number,
		DifficultyLevel: dbWord.DifficultyLevel,
		VerbConjugation: verbConj,
		Notes:           notes,
		CreatedAt:       dbWord.CreatedAt,
	}, nil
}

func (s *wordServiceImpl) ListWords(ctx context.Context, page, pageSize int) ([]models.Word, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	dbWords, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}

	// Convert DB models to API models
	apiWords := make([]models.Word, len(dbWords))
	for i, dbWord := range dbWords {
		var verbConj *string
		if len(dbWord.VerbConjugation) > 0 {
			s := string(dbWord.VerbConjugation)
			verbConj = &s
		}

		apiWords[i] = models.Word{
			ID:              dbWord.ID,
			Italian:         dbWord.Italian,
			English:         dbWord.English,
			PartsOfSpeech:   dbWord.PartsOfSpeech,
			Gender:          nil,
			Number:          nil,
			DifficultyLevel: dbWord.DifficultyLevel,
			VerbConjugation: verbConj,
			Notes:           nil,
			CreatedAt:       dbWord.CreatedAt,
		}

		if dbWord.Gender.Valid {
			apiWords[i].Gender = &dbWord.Gender.String
		}
		if dbWord.Number.Valid {
			apiWords[i].Number = &dbWord.Number.String
		}
		if dbWord.Notes.Valid {
			apiWords[i].Notes = &dbWord.Notes.String
		}
	}

	return apiWords, nil
}

func validateWord(word *models.Word) error {
	if word.Italian == "" || word.English == "" || word.PartsOfSpeech == "" {
		return ErrInvalid
	}
	if word.DifficultyLevel < 1 || word.DifficultyLevel > 5 {
		return ErrInvalid
	}
	return nil
}
