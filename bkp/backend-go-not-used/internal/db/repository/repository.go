package repository

import (
	"context"
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
	Search(ctx context.Context, query string, offset, limit int) ([]*models.Word, error)
	SearchCount(ctx context.Context, query string) (int, error)
	Count(ctx context.Context) (int, error)
}

type GroupRepository interface {
	Create(ctx context.Context, group *models.Group) error
	GetByID(ctx context.Context, id int64) (*models.Group, error)
	List(ctx context.Context, offset, limit int) ([]*models.Group, error)
	Update(ctx context.Context, group *models.Group) error
	Delete(ctx context.Context, id int64) error
	GetWordCount(ctx context.Context, groupID int64) (int, error)
	GetGroupStatistics(ctx context.Context, groupID int64) (apimodels.GroupStatistics, error)
	GetGroupProgress(ctx context.Context, groupID int64) (apimodels.GroupProgress, error)
}

type StudySessionRepository interface {
	Create(ctx context.Context, session *models.StudySession) error
	GetByID(ctx context.Context, id int64) (*models.StudySession, error)
	ListByGroupID(ctx context.Context, groupID int64, offset, limit int) ([]*models.StudySession, error)
	GetStats(ctx context.Context, groupID int64) (*models.StudyStats, error)
	List(ctx context.Context, page, limit int) ([]*models.StudySession, int, error)
}

type StudyActivityRepository interface {
	Create(ctx context.Context, activity *models.StudyActivity) error
	GetByID(ctx context.Context, id int64) (*models.StudyActivity, error)
	List(ctx context.Context, offset, limit int) ([]*models.StudyActivity, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetRecommended(ctx context.Context, count int) ([]*models.StudyActivity, error)
}

type WordReviewRepository interface {
	Create(ctx context.Context, review *models.WordReviewItem) error
	ListBySession(ctx context.Context, sessionID int64) ([]*models.WordReviewItem, error)
	GetWordStats(ctx context.Context, wordID int64) (*models.WordStats, error)
}

type SettingsRepository interface {
	GetSettings(ctx context.Context, userID int64) (*models.Settings, error)
	UpdateSettings(ctx context.Context, settings *models.Settings) error
	GetPreferences(ctx context.Context, userID int64) (*models.Preferences, error)
	UpdatePreferences(ctx context.Context, preferences *models.Preferences) error
}

type AnalyticsRepository interface {
	GetSessionAnalytics(ctx context.Context, userID int64) (*models.SessionAnalytics, error)
	GetSessionCalendar(ctx context.Context, userID int64) (*models.SessionCalendar, error)
}

type DashboardRepository interface {
	GetLastStudySession(ctx context.Context, userID int64) (*models.StudySession, error)
	GetStudyProgress(ctx context.Context, userID int64) (*models.StudyProgress, error)
	GetQuickStats(ctx context.Context, userID int64) (*models.QuickStats, error)
	GetStreak(ctx context.Context, userID int64) (*models.Streak, error)
	GetMasteryMetrics(ctx context.Context, userID int64) (*models.MasteryMetrics, error)
}


