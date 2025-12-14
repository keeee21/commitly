package mocks

import (
	"context"

	"github.com/keeee21/commitly/api/models"
)

// MockRivalRepository is a mock of IRivalRepository interface.
type MockRivalRepository struct {
	FindByUserIDFunc                       func(ctx context.Context, userID uint64) ([]models.Rival, error)
	FindByIDFunc                           func(ctx context.Context, id uint64) (*models.Rival, error)
	FindAllDistinctRivalsFunc              func(ctx context.Context) ([]models.Rival, error)
	CountByUserIDFunc                      func(ctx context.Context, userID uint64) (int64, error)
	CreateFunc                             func(ctx context.Context, rival *models.Rival) error
	DeleteFunc                             func(ctx context.Context, id uint64) error
	ExistsByUserIDAndRivalGithubUserIDFunc func(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error)
}

func (m *MockRivalRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.Rival, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockRivalRepository) FindByID(ctx context.Context, id uint64) (*models.Rival, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockRivalRepository) CountByUserID(ctx context.Context, userID uint64) (int64, error) {
	if m.CountByUserIDFunc != nil {
		return m.CountByUserIDFunc(ctx, userID)
	}
	return 0, nil
}

func (m *MockRivalRepository) Create(ctx context.Context, rival *models.Rival) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, rival)
	}
	return nil
}

func (m *MockRivalRepository) Delete(ctx context.Context, id uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockRivalRepository) ExistsByUserIDAndRivalGithubUserID(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error) {
	if m.ExistsByUserIDAndRivalGithubUserIDFunc != nil {
		return m.ExistsByUserIDAndRivalGithubUserIDFunc(ctx, userID, rivalGithubUserID)
	}
	return false, nil
}

func (m *MockRivalRepository) FindAllDistinctRivals(ctx context.Context) ([]models.Rival, error) {
	if m.FindAllDistinctRivalsFunc != nil {
		return m.FindAllDistinctRivalsFunc(ctx)
	}
	return nil, nil
}
