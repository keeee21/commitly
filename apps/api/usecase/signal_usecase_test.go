package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/keeee21/commitly/api/models"
	"github.com/stretchr/testify/assert"
)

type signalMockCircleRepository struct {
	FindByUserIDFunc       func(ctx context.Context, userID uint64) ([]models.Circle, error)
	FindByIDFunc           func(ctx context.Context, id uint64) (*models.Circle, error)
	FindByInviteCodeFunc   func(ctx context.Context, code string) (*models.Circle, error)
	CountByOwnerUserIDFunc func(ctx context.Context, userID uint64) (int64, error)
	CreateFunc             func(ctx context.Context, circle *models.Circle) error
	DeleteFunc             func(ctx context.Context, id uint64) error
	AddMemberFunc          func(ctx context.Context, member *models.CircleMember) error
	RemoveMemberFunc       func(ctx context.Context, circleID, userID uint64) error
	CountMembersFunc       func(ctx context.Context, circleID uint64) (int64, error)
	IsMemberFunc           func(ctx context.Context, circleID, userID uint64) (bool, error)
}

func (m *signalMockCircleRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.Circle, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID)
	}
	return nil, nil
}
func (m *signalMockCircleRepository) FindByID(ctx context.Context, id uint64) (*models.Circle, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *signalMockCircleRepository) FindByInviteCode(ctx context.Context, code string) (*models.Circle, error) {
	if m.FindByInviteCodeFunc != nil {
		return m.FindByInviteCodeFunc(ctx, code)
	}
	return nil, nil
}
func (m *signalMockCircleRepository) CountByOwnerUserID(ctx context.Context, userID uint64) (int64, error) {
	if m.CountByOwnerUserIDFunc != nil {
		return m.CountByOwnerUserIDFunc(ctx, userID)
	}
	return 0, nil
}
func (m *signalMockCircleRepository) Create(ctx context.Context, circle *models.Circle) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, circle)
	}
	return nil
}
func (m *signalMockCircleRepository) Delete(ctx context.Context, id uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
func (m *signalMockCircleRepository) AddMember(ctx context.Context, member *models.CircleMember) error {
	if m.AddMemberFunc != nil {
		return m.AddMemberFunc(ctx, member)
	}
	return nil
}
func (m *signalMockCircleRepository) RemoveMember(ctx context.Context, circleID, userID uint64) error {
	if m.RemoveMemberFunc != nil {
		return m.RemoveMemberFunc(ctx, circleID, userID)
	}
	return nil
}
func (m *signalMockCircleRepository) CountMembers(ctx context.Context, circleID uint64) (int64, error) {
	if m.CountMembersFunc != nil {
		return m.CountMembersFunc(ctx, circleID)
	}
	return 0, nil
}
func (m *signalMockCircleRepository) IsMember(ctx context.Context, circleID, userID uint64) (bool, error) {
	if m.IsMemberFunc != nil {
		return m.IsMemberFunc(ctx, circleID, userID)
	}
	return false, nil
}

type signalMockCommitStatsRepository struct {
	FindByGithubUserIDAndDateRangeFunc  func(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error)
	FindByGithubUserIDsAndDateRangeFunc func(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error)
	UpsertFunc                          func(ctx context.Context, stats *models.CommitStats) error
	UpsertBatchFunc                     func(ctx context.Context, statsList []models.CommitStats) error
}

func (m *signalMockCommitStatsRepository) FindByGithubUserIDAndDateRange(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	if m.FindByGithubUserIDAndDateRangeFunc != nil {
		return m.FindByGithubUserIDAndDateRangeFunc(ctx, githubUserID, startDate, endDate)
	}
	return nil, nil
}
func (m *signalMockCommitStatsRepository) FindByGithubUserIDsAndDateRange(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	if m.FindByGithubUserIDsAndDateRangeFunc != nil {
		return m.FindByGithubUserIDsAndDateRangeFunc(ctx, githubUserIDs, startDate, endDate)
	}
	return nil, nil
}
func (m *signalMockCommitStatsRepository) Upsert(ctx context.Context, stats *models.CommitStats) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(ctx, stats)
	}
	return nil
}
func (m *signalMockCommitStatsRepository) UpsertBatch(ctx context.Context, statsList []models.CommitStats) error {
	if m.UpsertBatchFunc != nil {
		return m.UpsertBatchFunc(ctx, statsList)
	}
	return nil
}

func makeCircleWithMembers() *models.Circle {
	return &models.Circle{
		ID:   1,
		Name: "テストサークル",
		Members: []models.CircleMember{
			{UserID: 1, User: models.User{ID: 1, GithubUserID: 100, GithubUsername: "me", AvatarURL: "https://avatar/me"}},
			{UserID: 2, User: models.User{ID: 2, GithubUserID: 200, GithubUsername: "tanaka", AvatarURL: "https://avatar/tanaka"}},
		},
	}
}

func TestGetSignals_SameDayCommit(t *testing.T) {
	circle := makeCircleWithMembers()
	today := time.Now().Truncate(24 * time.Hour)

	mockCircleRepo := &signalMockCircleRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return circle, nil
		},
	}
	mockCommitStatsRepo := &signalMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, ids []uint64, start, end time.Time) ([]models.CommitStats, error) {
			return []models.CommitStats{
				{GithubUserID: 100, Date: today, Repository: "me/repo", CommitCount: 3},
				{GithubUserID: 200, Date: today, Repository: "tanaka/repo", CommitCount: 2},
			}, nil
		},
	}

	uc := NewSignalUsecase(mockCircleRepo, mockCommitStatsRepo)
	signals, err := uc.GetSignals(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.NotEmpty(t, signals)

	hasSameDay := false
	for _, s := range signals {
		if s.Type == "same_day" {
			hasSameDay = true
			assert.Equal(t, "同じ日にコミット", s.Detail)
			assert.Contains(t, s.Usernames, "tanaka")
		}
	}
	assert.True(t, hasSameDay)
}

func TestGetSignals_SameHourCommit(t *testing.T) {
	circle := makeCircleWithMembers()
	today := time.Now().Truncate(24 * time.Hour)
	hour23 := 23

	mockCircleRepo := &signalMockCircleRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return circle, nil
		},
	}
	mockCommitStatsRepo := &signalMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, ids []uint64, start, end time.Time) ([]models.CommitStats, error) {
			return []models.CommitStats{
				{GithubUserID: 100, Date: today, Repository: "me/repo", CommitCount: 3, PrimaryHour: &hour23},
				{GithubUserID: 200, Date: today, Repository: "tanaka/repo", CommitCount: 2, PrimaryHour: &hour23},
			}, nil
		},
	}

	uc := NewSignalUsecase(mockCircleRepo, mockCommitStatsRepo)
	signals, err := uc.GetSignals(context.Background(), 1, 1)

	assert.NoError(t, err)

	hasSameHour := false
	for _, s := range signals {
		if s.Type == "same_hour" {
			hasSameHour = true
			assert.Equal(t, "23時台", s.Detail)
		}
	}
	assert.True(t, hasSameHour)
}

func TestGetSignals_SameLanguage(t *testing.T) {
	circle := makeCircleWithMembers()
	today := time.Now().Truncate(24 * time.Hour)

	mockCircleRepo := &signalMockCircleRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return circle, nil
		},
	}
	mockCommitStatsRepo := &signalMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, ids []uint64, start, end time.Time) ([]models.CommitStats, error) {
			return []models.CommitStats{
				{GithubUserID: 100, Date: today, Repository: "me/repo", CommitCount: 3, Language: "TypeScript"},
				{GithubUserID: 200, Date: today, Repository: "tanaka/repo", CommitCount: 2, Language: "TypeScript"},
			}, nil
		},
	}

	uc := NewSignalUsecase(mockCircleRepo, mockCommitStatsRepo)
	signals, err := uc.GetSignals(context.Background(), 1, 1)

	assert.NoError(t, err)

	hasSameLang := false
	for _, s := range signals {
		if s.Type == "same_language" {
			hasSameLang = true
			assert.Equal(t, "TypeScript", s.Detail)
		}
	}
	assert.True(t, hasSameLang)
}

func TestGetSignals_NoSignals(t *testing.T) {
	circle := makeCircleWithMembers()

	mockCircleRepo := &signalMockCircleRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return circle, nil
		},
	}
	mockCommitStatsRepo := &signalMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, ids []uint64, start, end time.Time) ([]models.CommitStats, error) {
			return []models.CommitStats{}, nil
		},
	}

	uc := NewSignalUsecase(mockCircleRepo, mockCommitStatsRepo)
	signals, err := uc.GetSignals(context.Background(), 1, 1)

	assert.NoError(t, err)
	assert.Empty(t, signals)
}

func TestGetSignals_NotMember(t *testing.T) {
	circle := makeCircleWithMembers()

	mockCircleRepo := &signalMockCircleRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return circle, nil
		},
	}
	mockCommitStatsRepo := &signalMockCommitStatsRepository{}

	uc := NewSignalUsecase(mockCircleRepo, mockCommitStatsRepo)
	_, err := uc.GetSignals(context.Background(), 999, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "メンバーではありません")
}

func TestGetSignals_CircleNotFound(t *testing.T) {
	mockCircleRepo := &signalMockCircleRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return nil, errors.New("not found")
		},
	}
	mockCommitStatsRepo := &signalMockCommitStatsRepository{}

	uc := NewSignalUsecase(mockCircleRepo, mockCommitStatsRepo)
	_, err := uc.GetSignals(context.Background(), 1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "サークルが見つかりません")
}

func TestGetRecentSignals_Success(t *testing.T) {
	today := time.Now().Truncate(24 * time.Hour)
	circle := makeCircleWithMembers()

	mockCircleRepo := &signalMockCircleRepository{
		FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.Circle, error) {
			return []models.Circle{*circle}, nil
		},
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return circle, nil
		},
	}
	mockCommitStatsRepo := &signalMockCommitStatsRepository{
		FindByGithubUserIDsAndDateRangeFunc: func(ctx context.Context, ids []uint64, start, end time.Time) ([]models.CommitStats, error) {
			return []models.CommitStats{
				{GithubUserID: 100, Date: today, Repository: "me/repo", CommitCount: 3},
				{GithubUserID: 200, Date: today, Repository: "tanaka/repo", CommitCount: 2},
			}, nil
		},
	}

	uc := NewSignalUsecase(mockCircleRepo, mockCommitStatsRepo)
	signals, err := uc.GetRecentSignals(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotEmpty(t, signals)
	assert.Equal(t, uint64(1), signals[0].CircleID)
	assert.Equal(t, "テストサークル", signals[0].CircleName)
}

func TestGetRecentSignals_NoCircles(t *testing.T) {
	mockCircleRepo := &signalMockCircleRepository{
		FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.Circle, error) {
			return []models.Circle{}, nil
		},
	}
	mockCommitStatsRepo := &signalMockCommitStatsRepository{}

	uc := NewSignalUsecase(mockCircleRepo, mockCommitStatsRepo)
	signals, err := uc.GetRecentSignals(context.Background(), 1)

	assert.NoError(t, err)
	assert.Empty(t, signals)
}
