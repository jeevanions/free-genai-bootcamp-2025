package services

import (
	"context"
	"errors"

	"github.com/jeevanions/italian-learning/internal/db/models"
	"github.com/jeevanions/italian-learning/internal/db/repository"
)

var (
	ErrInvalidReview = errors.New("invalid review data")
)

type WordReviewService interface {
	CreateReview(ctx context.Context, review *models.WordReviewItem) error
	ListSessionReviews(ctx context.Context, sessionID int64) ([]*models.WordReviewItem, error)
	GetWordStats(ctx context.Context, wordID int64) (*models.WordStats, error)
}

type wordReviewServiceImpl struct {
	reviewRepo  repository.WordReviewRepository
	wordRepo    repository.WordRepository
	sessionRepo repository.StudySessionRepository
}

func NewWordReviewService(
	reviewRepo repository.WordReviewRepository,
	wordRepo repository.WordRepository,
	sessionRepo repository.StudySessionRepository,
) WordReviewService {
	return &wordReviewServiceImpl{
		reviewRepo:  reviewRepo,
		wordRepo:    wordRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *wordReviewServiceImpl) CreateReview(ctx context.Context, review *models.WordReviewItem) error {
	if err := validateReview(review); err != nil {
		return err
	}

	// Verify word exists
	if _, err := s.wordRepo.GetByID(ctx, review.WordID); err != nil {
		return err
	}

	// Verify session exists
	if _, err := s.sessionRepo.GetByID(ctx, review.StudySessionID); err != nil {
		return err
	}

	return s.reviewRepo.Create(ctx, review)
}

func (s *wordReviewServiceImpl) ListSessionReviews(ctx context.Context, sessionID int64) ([]*models.WordReviewItem, error) {
	// Verify session exists
	if _, err := s.sessionRepo.GetByID(ctx, sessionID); err != nil {
		return nil, err
	}

	return s.reviewRepo.ListBySession(ctx, sessionID)
}

func (s *wordReviewServiceImpl) GetWordStats(ctx context.Context, wordID int64) (*models.WordStats, error) {
	// Verify word exists
	if _, err := s.wordRepo.GetByID(ctx, wordID); err != nil {
		return nil, err
	}

	stats, err := s.reviewRepo.GetWordStats(ctx, wordID)
	if err != nil {
		return nil, err
	}

	stats.CalculateAccuracy()
	return stats, nil
}

func validateReview(review *models.WordReviewItem) error {
	if review.WordID <= 0 || review.StudySessionID <= 0 {
		return ErrInvalidReview
	}
	return nil
}
