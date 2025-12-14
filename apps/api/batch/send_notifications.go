package batch

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/keeee21/commitly/api/gateway"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
)

// SendNotificationsConfig 通知送信バッチの設定
type SendNotificationsConfig struct {
	Period string // "weekly" or "monthly"
}

// ISendNotificationsDeps 通知送信バッチの依存関係インターフェース
type ISendNotificationsDeps interface {
	GetSlackNotificationRepo() repository.ISlackNotificationSettingRepository
	GetNotificationLogRepo() repository.INotificationLogRepository
	GetRivalRepo() repository.IRivalRepository
	GetCommitStatsRepo() repository.ICommitStatsRepository
	GetSlackGateway() gateway.ISlackGateway
}

// SendNotificationsDeps 通知送信バッチの依存関係
type SendNotificationsDeps struct {
	SlackNotificationRepo repository.ISlackNotificationSettingRepository
	NotificationLogRepo   repository.INotificationLogRepository
	RivalRepo             repository.IRivalRepository
	CommitStatsRepo       repository.ICommitStatsRepository
	SlackGateway          gateway.ISlackGateway
}

func (d *SendNotificationsDeps) GetSlackNotificationRepo() repository.ISlackNotificationSettingRepository {
	return d.SlackNotificationRepo
}

func (d *SendNotificationsDeps) GetNotificationLogRepo() repository.INotificationLogRepository {
	return d.NotificationLogRepo
}

func (d *SendNotificationsDeps) GetRivalRepo() repository.IRivalRepository {
	return d.RivalRepo
}

func (d *SendNotificationsDeps) GetCommitStatsRepo() repository.ICommitStatsRepository {
	return d.CommitStatsRepo
}

func (d *SendNotificationsDeps) GetSlackGateway() gateway.ISlackGateway {
	return d.SlackGateway
}

// DateRange 日付範囲
type DateRange struct {
	Start time.Time
	End   time.Time
}

// RunSendNotifications 通知送信バッチを実行
func RunSendNotifications(ctx context.Context, deps ISendNotificationsDeps, config SendNotificationsConfig) error {
	log.Println("Starting send-notifications batch...")
	startTime := time.Now()

	// 期間のバリデーション
	if config.Period != "weekly" && config.Period != "monthly" {
		return fmt.Errorf("invalid period: %s (must be 'weekly' or 'monthly')", config.Period)
	}

	log.Printf("Period: %s", config.Period)

	// 有効なSlack通知設定を取得
	slackSettings, err := deps.GetSlackNotificationRepo().FindAllEnabled(ctx)
	if err != nil {
		return fmt.Errorf("failed to get enabled Slack notification settings: %w", err)
	}

	log.Printf("Found %d enabled Slack notification settings", len(slackSettings))

	successCount := 0
	failCount := 0

	for _, setting := range slackSettings {
		var sendErr error
		var payload models.JSONPayload
		sentAt := time.Now()

		if config.Period == "weekly" {
			payload, sendErr = sendWeeklyReport(ctx, deps, setting)
		} else {
			payload, sendErr = sendMonthlyReport(ctx, deps, setting)
		}

		// ログを保存
		notificationLog := &models.NotificationLog{
			UserID:      setting.UserID,
			ChannelType: models.ChannelTypeSlack,
			Period:      config.Period,
			Payload:     payload,
			SentAt:      sentAt,
		}

		if sendErr != nil {
			notificationLog.Status = models.NotificationStatusFailed
			notificationLog.ErrorMessage = sendErr.Error()
			log.Printf("Failed to send notification for user %d: %v", setting.UserID, sendErr)
			failCount++
		} else {
			notificationLog.Status = models.NotificationStatusSuccess
			successCount++
		}

		// ログをDBに保存
		if err := deps.GetNotificationLogRepo().Create(ctx, notificationLog); err != nil {
			log.Printf("Failed to save notification log for user %d: %v", setting.UserID, err)
		}
	}

	elapsed := time.Since(startTime)
	log.Printf("send-notifications batch completed in %s (success: %d, failed: %d)", elapsed, successCount, failCount)

	return nil
}

// calculateWeeklyRange 週次レポートの日付範囲を計算（過去7日間）
func calculateWeeklyRange() DateRange {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	return DateRange{
		Start: today.AddDate(0, 0, -7),
		End:   today.AddDate(0, 0, -1),
	}
}

// calculateMonthlyRange 月次レポートの日付範囲を計算（先月）
func calculateMonthlyRange() (current DateRange, previous DateRange) {
	now := time.Now()

	// 先月の範囲
	firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	lastOfLastMonth := firstOfThisMonth.AddDate(0, 0, -1)
	firstOfLastMonth := time.Date(lastOfLastMonth.Year(), lastOfLastMonth.Month(), 1, 0, 0, 0, 0, time.Local)

	current = DateRange{
		Start: firstOfLastMonth,
		End:   lastOfLastMonth,
	}

	// 先々月の範囲
	lastOfTwoMonthsAgo := firstOfLastMonth.AddDate(0, 0, -1)
	firstOfTwoMonthsAgo := time.Date(lastOfTwoMonthsAgo.Year(), lastOfTwoMonthsAgo.Month(), 1, 0, 0, 0, 0, time.Local)

	previous = DateRange{
		Start: firstOfTwoMonthsAgo,
		End:   lastOfTwoMonthsAgo,
	}

	return current, previous
}

// sendWeeklyReport 週次レポートを送信
func sendWeeklyReport(
	ctx context.Context,
	deps ISendNotificationsDeps,
	setting models.SlackNotificationSetting,
) (models.JSONPayload, error) {
	user := setting.User
	dateRange := calculateWeeklyRange()

	// ユーザーのコミット統計を取得
	userStats, err := deps.GetCommitStatsRepo().FindByGithubUserIDAndDateRange(ctx, user.GithubUserID, dateRange.Start, dateRange.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get user commit stats: %w", err)
	}
	userCommits := sumCommits(userStats)

	// ライバルのコミット統計を取得
	rivalSummaries, err := getRivalSummaries(ctx, deps, setting.UserID, dateRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get rival summaries: %w", err)
	}

	// Slackメッセージを構築して送信
	message := gateway.BuildWeeklyReportMessage(
		user.GithubUsername,
		userCommits,
		rivalSummaries,
		dateRange.Start.Format("2006/01/02"),
		dateRange.End.Format("2006/01/02"),
	)

	// メッセージをJSONPayloadに変換
	payload := messageToPayload(message)

	if err := deps.GetSlackGateway().SendMessage(ctx, setting.WebhookURL, message); err != nil {
		return payload, fmt.Errorf("failed to send slack message: %w", err)
	}

	log.Printf("Sent weekly report to user %s (commits: %d)", user.GithubUsername, userCommits)
	return payload, nil
}

// sendMonthlyReport 月次レポートを送信
func sendMonthlyReport(
	ctx context.Context,
	deps ISendNotificationsDeps,
	setting models.SlackNotificationSetting,
) (models.JSONPayload, error) {
	user := setting.User
	currentRange, previousRange := calculateMonthlyRange()

	// 今月のコミット統計を取得
	currentStats, err := deps.GetCommitStatsRepo().FindByGithubUserIDAndDateRange(ctx, user.GithubUserID, currentRange.Start, currentRange.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get current month commit stats: %w", err)
	}
	currentCommits := sumCommits(currentStats)

	// 先月のコミット統計を取得
	previousStats, err := deps.GetCommitStatsRepo().FindByGithubUserIDAndDateRange(ctx, user.GithubUserID, previousRange.Start, previousRange.End)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous month commit stats: %w", err)
	}
	previousCommits := sumCommits(previousStats)

	// ライバルのコミット統計を取得（今月分）
	rivalSummaries, err := getRivalSummaries(ctx, deps, setting.UserID, currentRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get rival summaries: %w", err)
	}

	// 月のラベルを作成
	monthLabel := fmt.Sprintf("%d年%d月", currentRange.Start.Year(), currentRange.Start.Month())

	// Slackメッセージを構築して送信
	message := gateway.BuildMonthlyReportMessage(
		user.GithubUsername,
		gateway.MonthlyComparison{
			CurrentMonth:  currentCommits,
			PreviousMonth: previousCommits,
		},
		rivalSummaries,
		monthLabel,
	)

	// メッセージをJSONPayloadに変換
	payload := messageToPayload(message)

	if err := deps.GetSlackGateway().SendMessage(ctx, setting.WebhookURL, message); err != nil {
		return payload, fmt.Errorf("failed to send slack message: %w", err)
	}

	log.Printf("Sent monthly report to user %s (current: %d, previous: %d)", user.GithubUsername, currentCommits, previousCommits)
	return payload, nil
}

// getRivalSummaries ライバルのコミットサマリーを取得
func getRivalSummaries(
	ctx context.Context,
	deps ISendNotificationsDeps,
	userID uint64,
	dateRange DateRange,
) ([]gateway.RivalCommitSummary, error) {
	rivals, err := deps.GetRivalRepo().FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var summaries []gateway.RivalCommitSummary
	for _, rival := range rivals {
		rivalStats, err := deps.GetCommitStatsRepo().FindByGithubUserIDAndDateRange(ctx, rival.RivalGithubUserID, dateRange.Start, dateRange.End)
		if err != nil {
			log.Printf("Failed to get rival %s commit stats: %v", rival.RivalGithubUsername, err)
			continue
		}

		summaries = append(summaries, gateway.RivalCommitSummary{
			Username: rival.RivalGithubUsername,
			Commits:  sumCommits(rivalStats),
		})
	}

	return summaries, nil
}

func sumCommits(stats []models.CommitStats) int {
	total := 0
	for _, s := range stats {
		total += s.CommitCount
	}
	return total
}

// messageToPayload SlackMessageをJSONPayloadに変換
func messageToPayload(message *gateway.SlackMessage) models.JSONPayload {
	if message == nil {
		return nil
	}
	payload := models.JSONPayload{
		"blocks": message.Blocks,
	}
	if message.Text != "" {
		payload["text"] = message.Text
	}
	return payload
}
