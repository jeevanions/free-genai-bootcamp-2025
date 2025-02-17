package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jeevanions/italian-learning/internal/db/models"
)

// Error definitions
var (
	ErrNotFound = errors.New("record not found")
	ErrInvalid  = errors.New("invalid input")
)

// Interfaces
type WordRepository interface {
	Create(ctx context.Context, word *models.Word) error
	GetByID(ctx context.Context, id int64) (*models.Word, error)
	List(ctx context.Context, offset, limit int) ([]*models.Word, error)
	Update(ctx context.Context, word *models.Word) error
	Delete(ctx context.Context, id int64) error
}

type GroupRepository interface {
	Create(ctx context.Context, group *models.Group) error
	GetByID(ctx context.Context, id int64) (*models.Group, error)
	List(ctx context.Context, offset, limit int) ([]*models.Group, error)
	Update(ctx context.Context, group *models.Group) error
	Delete(ctx context.Context, id int64) error
}

type StudySessionRepository interface {
	Create(ctx context.Context, session *models.StudySession) error
	GetByID(ctx context.Context, id int64) (*models.StudySession, error)
	ListByGroupID(ctx context.Context, groupID int64, offset, limit int) ([]*models.StudySession, error)
	GetStats(ctx context.Context, groupID int64) (*models.StudyStats, error)
}

type StudyActivityRepository interface {
	Create(ctx context.Context, activity *models.StudyActivity) error
	GetByID(ctx context.Context, id int64) (*models.StudyActivity, error)
	List(ctx context.Context, offset, limit int) ([]*models.StudyActivity, error)
}

type WordReviewRepository interface {
	Create(ctx context.Context, review *models.WordReviewItem) error
	ListBySession(ctx context.Context, sessionID int64) ([]*models.WordReviewItem, error)
	GetWordStats(ctx context.Context, wordID int64) (*models.WordStats, error)
}

// Implementation for WordRepository
type SQLiteWordRepository struct {
	db *sql.DB
}

func NewSQLiteWordRepository(db *sql.DB) *SQLiteWordRepository {
	return &SQLiteWordRepository{db: db}
}

// Implement methods...
