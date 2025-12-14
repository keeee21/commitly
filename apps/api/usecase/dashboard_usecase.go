package usecase

import (
	"context"
	"time"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
)

// DailyCommitSummary 日別コミットサマリー
type DailyCommitSummary struct {
	Date        string `json:"date"`
	CommitCount int    `json:"commit_count"`
}

// RepositoryCommitSummary リポジトリ別コミットサマリー
type RepositoryCommitSummary struct {
	Repository  string `json:"repository"`
	CommitCount int    `json:"commit_count"`
}

// UserCommitStats ユーザーのコミット統計
type UserCommitStats struct {
	GithubUserID   uint64                    `json:"github_user_id"`
	GithubUsername string                    `json:"github_username"`
	AvatarURL      string                    `json:"avatar_url"`
	TotalCommits   int                       `json:"total_commits"`
	DailyStats     []DailyCommitSummary      `json:"daily_stats"`
	RepoStats      []RepositoryCommitSummary `json:"repo_stats"`
}

// DashboardData ダッシュボードデータ
type DashboardData struct {
	Period    string            `json:"period"` // "weekly" or "monthly"
	StartDate string            `json:"start_date"`
	EndDate   string            `json:"end_date"`
	MyStats   UserCommitStats   `json:"my_stats"`
	Rivals    []UserCommitStats `json:"rivals"`
}

// IDashboardUsecase ダッシュボードユースケースのインターフェース
type IDashboardUsecase interface {
	GetWeeklyDashboard(ctx context.Context, user *models.User, rivals []models.Rival) (*DashboardData, error)
	GetMonthlyDashboard(ctx context.Context, user *models.User, rivals []models.Rival) (*DashboardData, error)
}

type dashboardUsecase struct {
	commitStatsRepo repository.ICommitStatsRepository
}

// NewDashboardUsecase コンストラクタ
func NewDashboardUsecase(commitStatsRepo repository.ICommitStatsRepository) IDashboardUsecase {
	return &dashboardUsecase{
		commitStatsRepo: commitStatsRepo,
	}
}

func (u *dashboardUsecase) GetWeeklyDashboard(ctx context.Context, user *models.User, rivals []models.Rival) (*DashboardData, error) {
	now := time.Now()
	endDate := now
	startDate := now.AddDate(0, 0, -6) // 直近7日間

	return u.getDashboard(ctx, "weekly", startDate, endDate, user, rivals)
}

func (u *dashboardUsecase) GetMonthlyDashboard(ctx context.Context, user *models.User, rivals []models.Rival) (*DashboardData, error) {
	now := time.Now()
	endDate := now
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()) // 今月の1日

	return u.getDashboard(ctx, "monthly", startDate, endDate, user, rivals)
}

func (u *dashboardUsecase) getDashboard(ctx context.Context, period string, startDate, endDate time.Time, user *models.User, rivals []models.Rival) (*DashboardData, error) {
	// 対象のGithub User IDを収集
	githubUserIDs := []uint64{user.GithubUserID}
	for _, rival := range rivals {
		githubUserIDs = append(githubUserIDs, rival.RivalGithubUserID)
	}

	// コミット統計を取得
	stats, err := u.commitStatsRepo.FindByGithubUserIDsAndDateRange(ctx, githubUserIDs, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// ユーザー別に集計
	userStatsMap := make(map[uint64]*UserCommitStats)

	// 自分のデータを初期化
	userStatsMap[user.GithubUserID] = &UserCommitStats{
		GithubUserID:   user.GithubUserID,
		GithubUsername: user.GithubUsername,
		AvatarURL:      user.AvatarURL,
		DailyStats:     []DailyCommitSummary{},
		RepoStats:      []RepositoryCommitSummary{},
	}

	// ライバルのデータを初期化
	for _, rival := range rivals {
		userStatsMap[rival.RivalGithubUserID] = &UserCommitStats{
			GithubUserID:   rival.RivalGithubUserID,
			GithubUsername: rival.RivalGithubUsername,
			AvatarURL:      rival.RivalAvatarURL,
			DailyStats:     []DailyCommitSummary{},
			RepoStats:      []RepositoryCommitSummary{},
		}
	}

	// 日別・リポジトリ別に集計
	dailyMap := make(map[uint64]map[string]int)  // githubUserID -> date -> count
	repoMap := make(map[uint64]map[string]int)   // githubUserID -> repo -> count

	for _, stat := range stats {
		if dailyMap[stat.GithubUserID] == nil {
			dailyMap[stat.GithubUserID] = make(map[string]int)
		}
		if repoMap[stat.GithubUserID] == nil {
			repoMap[stat.GithubUserID] = make(map[string]int)
		}

		dateStr := stat.Date.Format("2006-01-02")
		dailyMap[stat.GithubUserID][dateStr] += stat.CommitCount
		repoMap[stat.GithubUserID][stat.Repository] += stat.CommitCount

		if userStats, ok := userStatsMap[stat.GithubUserID]; ok {
			userStats.TotalCommits += stat.CommitCount
		}
	}

	// 日別データを配列に変換
	for githubUserID, daily := range dailyMap {
		if userStats, ok := userStatsMap[githubUserID]; ok {
			for date, count := range daily {
				userStats.DailyStats = append(userStats.DailyStats, DailyCommitSummary{
					Date:        date,
					CommitCount: count,
				})
			}
		}
	}

	// リポジトリ別データを配列に変換
	for githubUserID, repos := range repoMap {
		if userStats, ok := userStatsMap[githubUserID]; ok {
			for repo, count := range repos {
				userStats.RepoStats = append(userStats.RepoStats, RepositoryCommitSummary{
					Repository:  repo,
					CommitCount: count,
				})
			}
		}
	}

	// ライバルの統計を配列に変換
	rivalStats := make([]UserCommitStats, 0, len(rivals))
	for _, rival := range rivals {
		if userStats, ok := userStatsMap[rival.RivalGithubUserID]; ok {
			rivalStats = append(rivalStats, *userStats)
		}
	}

	return &DashboardData{
		Period:    period,
		StartDate: startDate.Format("2006-01-02"),
		EndDate:   endDate.Format("2006-01-02"),
		MyStats:   *userStatsMap[user.GithubUserID],
		Rivals:    rivalStats,
	}, nil
}
