package mocks

import (
	"context"
	"time"

	"github.com/keeee21/commitly/api/models"
)

// MockCommitStatsRepository is a mock of ICommitStatsRepository interface.
type MockCommitStatsRepository struct {
	FindByGithubUserIDAndDateRangeFunc  func(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error)
	FindByGithubUserIDsAndDateRangeFunc func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error)
	UpsertFunc                          func(ctx context.Context, stats *models.CommitStats) error
	UpsertBatchFunc                     func(ctx context.Context, statsList []models.CommitStats) error
}

func (m *MockCommitStatsRepository) FindByGithubUserIDAndDateRange(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	if m.FindByGithubUserIDAndDateRangeFunc != nil {
		return m.FindByGithubUserIDAndDateRangeFunc(ctx, githubUserID, startDate, endDate)
	}
	return nil, nil
}

func (m *MockCommitStatsRepository) FindByGithubUserIDsAndDateRange(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	if m.FindByGithubUserIDsAndDateRangeFunc != nil {
		return m.FindByGithubUserIDsAndDateRangeFunc(ctx, githubUserIDs, startDate, endDate)
	}
	return nil, nil
}

func (m *MockCommitStatsRepository) Upsert(ctx context.Context, stats *models.CommitStats) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(ctx, stats)
	}
	return nil
}

func (m *MockCommitStatsRepository) UpsertBatch(ctx context.Context, statsList []models.CommitStats) error {
	if m.UpsertBatchFunc != nil {
		return m.UpsertBatchFunc(ctx, statsList)
	}
	return nil
}
