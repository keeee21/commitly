package batch

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/keeee21/commitly/api/gateway"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
	"github.com/stretchr/testify/assert"
)

// mockSlackNotificationSettingRepository テスト用のモック
type mockSlackNotificationSettingRepository struct {
	FindByUserIDFunc   func(ctx context.Context, userID uint64) (*models.SlackNotificationSetting, error)
	FindAllEnabledFunc func(ctx context.Context) ([]models.SlackNotificationSetting, error)
	UpsertFunc         func(ctx context.Context, setting *models.SlackNotificationSetting) error
	UpdateEnabledFunc  func(ctx context.Context, userID uint64, isEnabled bool) error
	DeleteFunc         func(ctx context.Context, userID uint64) error
}

func (m *mockSlackNotificationSettingRepository) FindByUserID(ctx context.Context, userID uint64) (*models.SlackNotificationSetting, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockSlackNotificationSettingRepository) FindAllEnabled(ctx context.Context) ([]models.SlackNotificationSetting, error) {
	if m.FindAllEnabledFunc != nil {
		return m.FindAllEnabledFunc(ctx)
	}
	return nil, nil
}

func (m *mockSlackNotificationSettingRepository) Upsert(ctx context.Context, setting *models.SlackNotificationSetting) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(ctx, setting)
	}
	return nil
}

func (m *mockSlackNotificationSettingRepository) UpdateEnabled(ctx context.Context, userID uint64, isEnabled bool) error {
	if m.UpdateEnabledFunc != nil {
		return m.UpdateEnabledFunc(ctx, userID, isEnabled)
	}
	return nil
}

func (m *mockSlackNotificationSettingRepository) Delete(ctx context.Context, userID uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, userID)
	}
	return nil
}

// mockRivalRepository テスト用のモック
type mockRivalRepository struct {
	FindByUserIDFunc                       func(ctx context.Context, userID uint64) ([]models.Rival, error)
	FindByIDFunc                           func(ctx context.Context, id uint64) (*models.Rival, error)
	FindAllDistinctRivalsFunc              func(ctx context.Context) ([]models.Rival, error)
	CountByUserIDFunc                      func(ctx context.Context, userID uint64) (int64, error)
	CreateFunc                             func(ctx context.Context, rival *models.Rival) error
	DeleteFunc                             func(ctx context.Context, id uint64) error
	ExistsByUserIDAndRivalGithubUserIDFunc func(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error)
}

func (m *mockRivalRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.Rival, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockRivalRepository) FindByID(ctx context.Context, id uint64) (*models.Rival, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRivalRepository) FindAllDistinctRivals(ctx context.Context) ([]models.Rival, error) {
	if m.FindAllDistinctRivalsFunc != nil {
		return m.FindAllDistinctRivalsFunc(ctx)
	}
	return nil, nil
}

func (m *mockRivalRepository) CountByUserID(ctx context.Context, userID uint64) (int64, error) {
	if m.CountByUserIDFunc != nil {
		return m.CountByUserIDFunc(ctx, userID)
	}
	return 0, nil
}

func (m *mockRivalRepository) Create(ctx context.Context, rival *models.Rival) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, rival)
	}
	return nil
}

func (m *mockRivalRepository) Delete(ctx context.Context, id uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *mockRivalRepository) ExistsByUserIDAndRivalGithubUserID(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error) {
	if m.ExistsByUserIDAndRivalGithubUserIDFunc != nil {
		return m.ExistsByUserIDAndRivalGithubUserIDFunc(ctx, userID, rivalGithubUserID)
	}
	return false, nil
}

// mockCommitStatsRepository テスト用のモック
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

// mockSlackGateway テスト用のモック
type mockSlackGateway struct {
	SendMessageFunc func(ctx context.Context, webhookURL string, message *gateway.SlackMessage) error
}

func (m *mockSlackGateway) SendMessage(ctx context.Context, webhookURL string, message *gateway.SlackMessage) error {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(ctx, webhookURL, message)
	}
	return nil
}

// mockNotificationLogRepository テスト用のモック
type mockNotificationLogRepository struct {
	CreateFunc              func(ctx context.Context, log *models.NotificationLog) error
	FindByUserIDFunc        func(ctx context.Context, userID uint64, limit int) ([]models.NotificationLog, error)
	FindByDateRangeFunc     func(ctx context.Context, startDate, endDate time.Time) ([]models.NotificationLog, error)
	FindByPeriodAndDateFunc func(ctx context.Context, period string, date time.Time) ([]models.NotificationLog, error)
}

func (m *mockNotificationLogRepository) Create(ctx context.Context, log *models.NotificationLog) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, log)
	}
	return nil
}

func (m *mockNotificationLogRepository) FindByUserID(ctx context.Context, userID uint64, limit int) ([]models.NotificationLog, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID, limit)
	}
	return nil, nil
}

func (m *mockNotificationLogRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NotificationLog, error) {
	if m.FindByDateRangeFunc != nil {
		return m.FindByDateRangeFunc(ctx, startDate, endDate)
	}
	return nil, nil
}

func (m *mockNotificationLogRepository) FindByPeriodAndDate(ctx context.Context, period string, date time.Time) ([]models.NotificationLog, error) {
	if m.FindByPeriodAndDateFunc != nil {
		return m.FindByPeriodAndDateFunc(ctx, period, date)
	}
	return nil, nil
}

// testDeps テスト用の依存関係
type testDeps struct {
	slackNotificationRepo *mockSlackNotificationSettingRepository
	notificationLogRepo   *mockNotificationLogRepository
	rivalRepo             *mockRivalRepository
	commitStatsRepo       *mockCommitStatsRepository
	slackGateway          *mockSlackGateway
}

func (d *testDeps) GetSlackNotificationRepo() repository.ISlackNotificationSettingRepository {
	return d.slackNotificationRepo
}

func (d *testDeps) GetNotificationLogRepo() repository.INotificationLogRepository {
	return d.notificationLogRepo
}

func (d *testDeps) GetRivalRepo() repository.IRivalRepository {
	return d.rivalRepo
}

func (d *testDeps) GetCommitStatsRepo() repository.ICommitStatsRepository {
	return d.commitStatsRepo
}

func (d *testDeps) GetSlackGateway() gateway.ISlackGateway {
	return d.slackGateway
}

func TestCalculateWeeklyRange(t *testing.T) {
	dateRange := calculateWeeklyRange()

	assert.True(t, dateRange.Start.Before(dateRange.End))
	assert.Equal(t, 6, int(dateRange.End.Sub(dateRange.Start).Hours()/24)) // 7日間の差
}

func TestCalculateMonthlyRange(t *testing.T) {
	current, previous := calculateMonthlyRange()

	// 今月の範囲は先月を指す
	assert.True(t, current.Start.Before(current.End) || current.Start.Equal(current.End))
	// 先月の範囲は先々月を指す
	assert.True(t, previous.Start.Before(previous.End) || previous.Start.Equal(previous.End))
	// 先月は先々月より後
	assert.True(t, current.Start.After(previous.End))
}

func TestSumCommits(t *testing.T) {
	stats := []models.CommitStats{
		{CommitCount: 5},
		{CommitCount: 3},
		{CommitCount: 2},
	}

	total := sumCommits(stats)

	assert.Equal(t, 10, total)
}

func TestSumCommits_Empty(t *testing.T) {
	stats := []models.CommitStats{}

	total := sumCommits(stats)

	assert.Equal(t, 0, total)
}

func TestRunSendNotifications_InvalidPeriod(t *testing.T) {
	ctx := context.Background()

	deps := &testDeps{
		slackNotificationRepo: &mockSlackNotificationSettingRepository{},
		notificationLogRepo:   &mockNotificationLogRepository{},
		rivalRepo:             &mockRivalRepository{},
		commitStatsRepo:       &mockCommitStatsRepository{},
		slackGateway:          &mockSlackGateway{},
	}

	config := SendNotificationsConfig{Period: "daily"} // invalid
	err := RunSendNotifications(ctx, deps, config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid period")
}

func TestRunSendNotifications_NoEnabledSettings(t *testing.T) {
	ctx := context.Background()

	deps := &testDeps{
		slackNotificationRepo: &mockSlackNotificationSettingRepository{
			FindAllEnabledFunc: func(ctx context.Context) ([]models.SlackNotificationSetting, error) {
				return []models.SlackNotificationSetting{}, nil
			},
		},
		notificationLogRepo: &mockNotificationLogRepository{},
		rivalRepo:           &mockRivalRepository{},
		commitStatsRepo:     &mockCommitStatsRepository{},
		slackGateway:        &mockSlackGateway{},
	}

	config := SendNotificationsConfig{Period: "weekly"}
	err := RunSendNotifications(ctx, deps, config)

	assert.NoError(t, err)
}

func TestRunSendNotifications_RepositoryError(t *testing.T) {
	ctx := context.Background()

	deps := &testDeps{
		slackNotificationRepo: &mockSlackNotificationSettingRepository{
			FindAllEnabledFunc: func(ctx context.Context) ([]models.SlackNotificationSetting, error) {
				return nil, errors.New("database error")
			},
		},
		notificationLogRepo: &mockNotificationLogRepository{},
		rivalRepo:           &mockRivalRepository{},
		commitStatsRepo:     &mockCommitStatsRepository{},
		slackGateway:        &mockSlackGateway{},
	}

	config := SendNotificationsConfig{Period: "weekly"}
	err := RunSendNotifications(ctx, deps, config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get enabled Slack notification settings")
}

func TestRunSendNotifications_WeeklySuccess(t *testing.T) {
	ctx := context.Background()
	slackCalled := false
	logSaved := false

	deps := &testDeps{
		slackNotificationRepo: &mockSlackNotificationSettingRepository{
			FindAllEnabledFunc: func(ctx context.Context) ([]models.SlackNotificationSetting, error) {
				return []models.SlackNotificationSetting{
					{
						ID:         1,
						UserID:     1,
						WebhookURL: "https://hooks.slack.com/test",
						IsEnabled:  true,
						User: models.User{
							ID:             1,
							GithubUserID:   12345,
							GithubUsername: "testuser",
						},
					},
				}, nil
			},
		},
		notificationLogRepo: &mockNotificationLogRepository{
			CreateFunc: func(ctx context.Context, log *models.NotificationLog) error {
				logSaved = true
				assert.Equal(t, models.NotificationStatusSuccess, log.Status)
				assert.Equal(t, "weekly", log.Period)
				assert.NotNil(t, log.Payload)
				assert.Contains(t, log.Payload, "blocks")
				return nil
			},
		},
		rivalRepo: &mockRivalRepository{
			FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
				return []models.Rival{
					{
						RivalGithubUserID:   67890,
						RivalGithubUsername: "rival1",
					},
				}, nil
			},
		},
		commitStatsRepo: &mockCommitStatsRepository{
			FindByGithubUserIDAndDateRangeFunc: func(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
				if githubUserID == 12345 {
					return []models.CommitStats{{CommitCount: 10}}, nil
				}
				return []models.CommitStats{{CommitCount: 5}}, nil
			},
		},
		slackGateway: &mockSlackGateway{
			SendMessageFunc: func(ctx context.Context, webhookURL string, message *gateway.SlackMessage) error {
				slackCalled = true
				assert.Equal(t, "https://hooks.slack.com/test", webhookURL)
				return nil
			},
		},
	}

	config := SendNotificationsConfig{Period: "weekly"}
	err := RunSendNotifications(ctx, deps, config)

	assert.NoError(t, err)
	assert.True(t, slackCalled)
	assert.True(t, logSaved)
}

func TestRunSendNotifications_MonthlySuccess(t *testing.T) {
	ctx := context.Background()
	slackCalled := false
	logSaved := false

	deps := &testDeps{
		slackNotificationRepo: &mockSlackNotificationSettingRepository{
			FindAllEnabledFunc: func(ctx context.Context) ([]models.SlackNotificationSetting, error) {
				return []models.SlackNotificationSetting{
					{
						ID:         1,
						UserID:     1,
						WebhookURL: "https://hooks.slack.com/test",
						IsEnabled:  true,
						User: models.User{
							ID:             1,
							GithubUserID:   12345,
							GithubUsername: "testuser",
						},
					},
				}, nil
			},
		},
		notificationLogRepo: &mockNotificationLogRepository{
			CreateFunc: func(ctx context.Context, log *models.NotificationLog) error {
				logSaved = true
				assert.Equal(t, models.NotificationStatusSuccess, log.Status)
				assert.Equal(t, "monthly", log.Period)
				assert.NotNil(t, log.Payload)
				assert.Contains(t, log.Payload, "blocks")
				return nil
			},
		},
		rivalRepo: &mockRivalRepository{
			FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
				return []models.Rival{
					{
						RivalGithubUserID:   67890,
						RivalGithubUsername: "rival1",
					},
				}, nil
			},
		},
		commitStatsRepo: &mockCommitStatsRepository{
			FindByGithubUserIDAndDateRangeFunc: func(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
				if githubUserID == 12345 {
					return []models.CommitStats{{CommitCount: 20}}, nil
				}
				return []models.CommitStats{{CommitCount: 15}}, nil
			},
		},
		slackGateway: &mockSlackGateway{
			SendMessageFunc: func(ctx context.Context, webhookURL string, message *gateway.SlackMessage) error {
				slackCalled = true
				assert.Equal(t, "https://hooks.slack.com/test", webhookURL)
				return nil
			},
		},
	}

	config := SendNotificationsConfig{Period: "monthly"}
	err := RunSendNotifications(ctx, deps, config)

	assert.NoError(t, err)
	assert.True(t, slackCalled)
	assert.True(t, logSaved)
}
