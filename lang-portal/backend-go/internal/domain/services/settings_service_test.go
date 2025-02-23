package services

import (
	"errors"
	"testing"

	"github.com/jeevanions/lang-portal/backend-go/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSeeder is a mock implementation of the seeder.Seeder interface
type MockSeeder struct {
	mock.Mock
}

func (m *MockSeeder) SeedFromJSON(seedDir string) error {
	args := m.Called(seedDir)
	return args.Error(0)
}

func TestSettingsService_ResetHistory(t *testing.T) {
	t.Run("successful reset", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockSeeder := new(MockSeeder)
		service := NewSettingsService(mockRepo, mockSeeder)

		mockRepo.On("ResetHistory").Return(nil)
		err := service.ResetHistory()
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockSeeder := new(MockSeeder)
		service := NewSettingsService(mockRepo, mockSeeder)

		expectedErr := errors.New("database error")
		mockRepo.On("ResetHistory").Return(expectedErr)

		err := service.ResetHistory()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSettingsService_FullReset(t *testing.T) {
	t.Run("successful reset", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockSeeder := new(MockSeeder)
		service := NewSettingsService(mockRepo, mockSeeder)

		mockRepo.On("DropAllTables").Return(nil)
		mockRepo.On("CreateTables").Return(nil)
		mockSeeder.On("SeedFromJSON", "internal/db/seeds").Return(nil)

		err := service.FullReset()

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockSeeder.AssertExpectations(t)
	})

	t.Run("drop tables error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockSeeder := new(MockSeeder)
		service := NewSettingsService(mockRepo, mockSeeder)

		mockRepo.On("DropAllTables").Return(errors.New("drop error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "drop error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("create tables error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockSeeder := new(MockSeeder)
		service := NewSettingsService(mockRepo, mockSeeder)

		mockRepo.On("DropAllTables").Return(nil)
		mockRepo.On("CreateTables").Return(errors.New("create error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "create error", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("seed error", func(t *testing.T) {
		mockRepo := new(mocks.MockRepository)
		mockSeeder := new(MockSeeder)
		service := NewSettingsService(mockRepo, mockSeeder)

		mockRepo.On("DropAllTables").Return(nil)
		mockRepo.On("CreateTables").Return(nil)
		mockSeeder.On("SeedFromJSON", "internal/db/seeds").Return(errors.New("seed error"))

		err := service.FullReset()

		assert.Error(t, err)
		assert.Equal(t, "seed error", err.Error())
		mockRepo.AssertExpectations(t)
		mockSeeder.AssertExpectations(t)
	})
}
