package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/keeee21/commitly/api/models"
	"github.com/stretchr/testify/assert"
)

// mockCommitStatsRepository テスト用のモックリポジトリ
type mockCommitStatsRepository struct {
	FindByGithubUserIDAndDateRangeFunc  func(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error)
	FindByGithubUserIDsAndDateRangeFunc func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error)
	UpsertFunc                          func(ctx context.Context, stats *models.CommitStats) error
	UpsertBatchFunc                     func(ctx context.Context, statsList []models.CommitStats) error
}

func (m *mockCommitStatsRepository) FindByGithubUserIDAndDateRange(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	if m.FindByGithubUserIDAndDateRangeFunc != nil {
		return m.FindByGithubUserIDAndDateRangeFunc(ctx, githubUserID, startDate, endDate)
	}
	return nil, nil
}

func (m *mockCommitStatsRepository) FindByGithubUserIDsAndDateRange(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	if m.FindByGithubUserIDsAndDateRangeFunc != nil {
		return m.FindByGithubUserIDsAndDateRangeFunc(ctx, githubUserIDs, startDate, endDate)
	}
	return nil, nil
}

func (m *mockCommitStatsRepository) Upsert(ctx context.Context, stats *models.CommitStats) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(ctx, stats)
	}
	return nil
}

func (m *mockCommitStatsRepository) UpsertBatch(ctx context.Context, statsList []models.CommitStats) error {
	if m.UpsertBatchFunc != nil {
		return m.UpsertBatchFunc(ctx, statsList)
	}
	return nil
}

func TestGetWeeklyDashboard_Success(t *testing.T) {
	ctx := context.Background()

	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
		AvatarURL:      "https://avatar.url/testuser",
	}

	rivals := []models.Rival{
		{
			ID:                  1,
			UserID:              1,
			RivalGithubUserID:   200,
			RivalGithubUsername: "rival1",
			RivalAvatarURL:      "https://avatar.url/rival1",
		},
	}

	now := time.Now()
	mockStats := []models.CommitStats{
		{
			GithubUserID:   100,
			GithubUsername: "testuser",
			Date:           now,
			Repository:     "testuser/repo1",
			CommitCount:    5,
		},
		{
			GithubUserID:   100,
			GithubUsername: "testuser",
			Date:           now.AddDate(0, 0, -1),
			Repository:     "testuser/repo2",
			CommitCount:    3,
		},
		{
			GithubUserID:   200,
			GithubUsername: "rival1",
			Date:           now,
			Repository:     "rival1/repo1",
			CommitCount:    10,
		},
	}

	mockRepo := &mockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			return mockStats, nil
		},
	}

	usecase := NewDashboardUsecase(mockRepo)
	data, err := usecase.GetWeeklyDashboard(ctx, user, rivals)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "weekly", data.Period)
	assert.Equal(t, 8, data.MyStats.TotalCommits)
	assert.Equal(t, "testuser", data.MyStats.GithubUsername)
	assert.Len(t, data.Rivals, 1)
	assert.Equal(t, 10, data.Rivals[0].TotalCommits)
	assert.Equal(t, "rival1", data.Rivals[0].GithubUsername)
}

func TestGetWeeklyDashboard_NoCommits(t *testing.T) {
	ctx := context.Background()

	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
	}

	rivals := []models.Rival{}

	mockRepo := &mockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			return []models.CommitStats{}, nil
		},
	}

	usecase := NewDashboardUsecase(mockRepo)
	data, err := usecase.GetWeeklyDashboard(ctx, user, rivals)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 0, data.MyStats.TotalCommits)
	assert.Empty(t, data.Rivals)
}

func TestGetWeeklyDashboard_RepositoryError(t *testing.T) {
	ctx := context.Background()

	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
	}

	rivals := []models.Rival{}

	mockRepo := &mockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			return nil, errors.New("database error")
		},
	}

	usecase := NewDashboardUsecase(mockRepo)
	data, err := usecase.GetWeeklyDashboard(ctx, user, rivals)

	assert.Error(t, err)
	assert.Nil(t, data)
	assert.Equal(t, "database error", err.Error())
}

func TestGetMonthlyDashboard_Success(t *testing.T) {
	ctx := context.Background()

	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
	}

	rivals := []models.Rival{}

	now := time.Now()
	mockStats := []models.CommitStats{
		{
			GithubUserID:   100,
			GithubUsername: "testuser",
			Date:           now,
			Repository:     "testuser/repo1",
			CommitCount:    15,
		},
	}

	mockRepo := &mockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			// 今月の1日から開始していることを確認
			expectedStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			assert.Equal(t, expectedStart.Format("2006-01-02"), startDate.Format("2006-01-02"))
			return mockStats, nil
		},
	}

	usecase := NewDashboardUsecase(mockRepo)
	data, err := usecase.GetMonthlyDashboard(ctx, user, rivals)

	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "monthly", data.Period)
	assert.Equal(t, 15, data.MyStats.TotalCommits)
}

func TestGetDashboard_MultipleRivals(t *testing.T) {
	ctx := context.Background()

	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
	}

	rivals := []models.Rival{
		{RivalGithubUserID: 200, RivalGithubUsername: "rival1"},
		{RivalGithubUserID: 300, RivalGithubUsername: "rival2"},
		{RivalGithubUserID: 400, RivalGithubUsername: "rival3"},
	}

	now := time.Now()
	mockStats := []models.CommitStats{
		{GithubUserID: 100, Date: now, Repository: "repo", CommitCount: 5},
		{GithubUserID: 200, Date: now, Repository: "repo", CommitCount: 10},
		{GithubUserID: 300, Date: now, Repository: "repo", CommitCount: 15},
		{GithubUserID: 400, Date: now, Repository: "repo", CommitCount: 20},
	}

	mockRepo := &mockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			assert.Len(t, githubUserIDs, 4) // user + 3 rivals
			return mockStats, nil
		},
	}

	usecase := NewDashboardUsecase(mockRepo)
	data, err := usecase.GetWeeklyDashboard(ctx, user, rivals)

	assert.NoError(t, err)
	assert.Len(t, data.Rivals, 3)
}
