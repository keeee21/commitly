package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/keeee21/commitly/api/models"
	"github.com/stretchr/testify/assert"
)

// activityMockCommitStatsRepository テスト用のモックリポジトリ
type activityMockCommitStatsRepository struct {
	FindByGithubUserIDAndDateRangeFunc  func(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error)
	FindByGithubUserIDsAndDateRangeFunc func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error)
	UpsertFunc                          func(ctx context.Context, stats *models.CommitStats) error
	UpsertBatchFunc                     func(ctx context.Context, statsList []models.CommitStats) error
}

func (m *activityMockCommitStatsRepository) FindByGithubUserIDAndDateRange(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	if m.FindByGithubUserIDAndDateRangeFunc != nil {
		return m.FindByGithubUserIDAndDateRangeFunc(ctx, githubUserID, startDate, endDate)
	}
	return nil, nil
}

func (m *activityMockCommitStatsRepository) FindByGithubUserIDsAndDateRange(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	if m.FindByGithubUserIDsAndDateRangeFunc != nil {
		return m.FindByGithubUserIDsAndDateRangeFunc(ctx, githubUserIDs, startDate, endDate)
	}
	return nil, nil
}

func (m *activityMockCommitStatsRepository) Upsert(ctx context.Context, stats *models.CommitStats) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(ctx, stats)
	}
	return nil
}

func (m *activityMockCommitStatsRepository) UpsertBatch(ctx context.Context, statsList []models.CommitStats) error {
	if m.UpsertBatchFunc != nil {
		return m.UpsertBatchFunc(ctx, statsList)
	}
	return nil
}

func TestGetActivityStream_Success(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
		AvatarURL:      "https://avatar.url/1",
	}
	rivals := []models.Rival{
		{
			ID:                  1,
			RivalGithubUserID:   200,
			RivalGithubUsername: "rival1",
			RivalAvatarURL:      "https://avatar.url/2",
		},
	}

	mockRepo := &activityMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			return []models.CommitStats{
				{GithubUserID: 100, GithubUsername: "testuser", Date: now, Repository: "my-repo", CommitCount: 3},
				{GithubUserID: 200, GithubUsername: "rival1", Date: now, Repository: "rival-repo", CommitCount: 5},
			}, nil
		},
	}

	uc := NewActivityUsecase(mockRepo)
	result, err := uc.GetActivityStream(ctx, user, rivals)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Activities, 2)
}

func TestGetActivityStream_NoRivals(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
		AvatarURL:      "https://avatar.url/1",
	}

	mockRepo := &activityMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			return []models.CommitStats{
				{GithubUserID: 100, GithubUsername: "testuser", Date: now, Repository: "my-repo", CommitCount: 2},
			}, nil
		},
	}

	uc := NewActivityUsecase(mockRepo)
	result, err := uc.GetActivityStream(ctx, user, []models.Rival{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Activities, 1)
	assert.Equal(t, "testuser", result.Activities[0].GithubUsername)
}

func TestGetActivityStream_Empty(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
	}

	mockRepo := &activityMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			return []models.CommitStats{}, nil
		},
	}

	uc := NewActivityUsecase(mockRepo)
	result, err := uc.GetActivityStream(ctx, user, []models.Rival{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result.Activities)
}

func TestGetActivityStream_Error(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
	}

	mockRepo := &activityMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			return nil, errors.New("database error")
		},
	}

	uc := NewActivityUsecase(mockRepo)
	result, err := uc.GetActivityStream(ctx, user, []models.Rival{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "アクティビティデータの取得に失敗しました")
}

func TestGetRhythm_Success(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
		AvatarURL:      "https://avatar.url/1",
	}
	rivals := []models.Rival{
		{
			RivalGithubUserID:   200,
			RivalGithubUsername: "rival1",
			RivalAvatarURL:      "https://avatar.url/2",
		},
	}

	// 安定型: 月〜金の5日
	now := time.Now()
	mockRepo := &activityMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			stats := []models.CommitStats{}
			// 直近7日分のデータを生成
			for i := 0; i < 7; i++ {
				d := now.AddDate(0, 0, -i)
				stats = append(stats, models.CommitStats{
					GithubUserID:   100,
					GithubUsername: "testuser",
					Date:           d,
					Repository:     "repo",
					CommitCount:    1,
				})
			}
			return stats, nil
		},
	}

	uc := NewActivityUsecase(mockRepo)
	result, err := uc.GetRhythm(ctx, user, rivals)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 2)
	assert.Equal(t, "testuser", result.Users[0].GithubUsername)
	assert.Equal(t, "安定型", result.Users[0].PatternLabel)
}

func TestGetRhythm_WeekendPattern(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
		AvatarURL:      "https://avatar.url/1",
	}

	mockRepo := &activityMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			// 土日のみのデータ
			stats := []models.CommitStats{}
			now := time.Now()
			for i := 0; i < 7; i++ {
				d := now.AddDate(0, 0, -i)
				if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
					stats = append(stats, models.CommitStats{
						GithubUserID:   100,
						GithubUsername: "testuser",
						Date:           d,
						Repository:     "repo",
						CommitCount:    3,
					})
				}
			}
			return stats, nil
		},
	}

	uc := NewActivityUsecase(mockRepo)
	result, err := uc.GetRhythm(ctx, user, []models.Rival{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "週末型", result.Users[0].PatternLabel)
}

func TestGetRhythm_BurstPattern(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
		AvatarURL:      "https://avatar.url/1",
	}

	mockRepo := &activityMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			// 平日3日のデータ（安定型でも週末型でもない）
			stats := []models.CommitStats{}
			now := time.Now()
			count := 0
			for i := 0; i < 7 && count < 3; i++ {
				d := now.AddDate(0, 0, -i)
				if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
					stats = append(stats, models.CommitStats{
						GithubUserID:   100,
						GithubUsername: "testuser",
						Date:           d,
						Repository:     "repo",
						CommitCount:    5,
					})
					count++
				}
			}
			return stats, nil
		},
	}

	uc := NewActivityUsecase(mockRepo)
	result, err := uc.GetRhythm(ctx, user, []models.Rival{})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "バースト型", result.Users[0].PatternLabel)
}

func TestGetRhythm_Error(t *testing.T) {
	ctx := context.Background()
	user := &models.User{
		ID:             1,
		GithubUserID:   100,
		GithubUsername: "testuser",
	}

	mockRepo := &activityMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
			return nil, errors.New("database error")
		},
	}

	uc := NewActivityUsecase(mockRepo)
	result, err := uc.GetRhythm(ctx, user, []models.Rival{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "リズムデータの取得に失敗しました")
}

func TestClassifyPattern(t *testing.T) {
	assert.Equal(t, "安定型", classifyPattern(5, 4, 1))
	assert.Equal(t, "安定型", classifyPattern(7, 5, 2))
	assert.Equal(t, "週末型", classifyPattern(2, 0, 2))
	assert.Equal(t, "週末型", classifyPattern(3, 2, 1))
	assert.Equal(t, "バースト型", classifyPattern(3, 3, 0))
	assert.Equal(t, "バースト型", classifyPattern(4, 4, 0))
}
