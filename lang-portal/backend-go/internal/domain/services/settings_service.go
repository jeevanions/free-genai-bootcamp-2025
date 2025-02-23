package services

import (
	"github.com/jeevanions/lang-portal/backend-go/internal/db/repository"
)

type Seeder interface {
	SeedFromJSON(seedDir string) error
}

type SettingsServiceInterface interface {
	ResetHistory() error
	FullReset() error
}

type SettingsService struct {
	repo   repository.Repository
	seeder Seeder
}

func NewSettingsService(repo repository.Repository, seeder Seeder) *SettingsService {
	return &SettingsService{
		repo:   repo,
		seeder: seeder,
	}
}

func (s *SettingsService) ResetHistory() error {
	return s.repo.ResetHistory()
}

func (s *SettingsService) FullReset() error {
	// First drop all tables
	if err := s.repo.DropAllTables(); err != nil {
		return err
	}

	// Then recreate tables and seed data
	if err := s.repo.CreateTables(); err != nil {
		return err
	}

	return s.seeder.SeedFromJSON("internal/db/seeds")
}


