package usecase

import (
	"context"
	"log"
	"time"

	"github.com/keeee21/commitly/api/gateway"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
)

// ISyncCommitsUsecase コミット同期ユースケースのインターフェース
type ISyncCommitsUsecase interface {
	SyncAllUsers(ctx context.Context) error
	SyncAllUsersWithDateRange(ctx context.Context, fromDate, toDate *time.Time) error
	SyncUser(ctx context.Context, githubUserID uint64, githubUsername string, fromDate, toDate *time.Time) error
}

type syncCommitsUsecase struct {
	userRepo        repository.IUserRepository
	rivalRepo       repository.IRivalRepository
	commitStatsRepo repository.ICommitStatsRepository
	githubGateway   gateway.IGithubGateway
}

// NewSyncCommitsUsecase コンストラクタ
func NewSyncCommitsUsecase(
	userRepo repository.IUserRepository,
	rivalRepo repository.IRivalRepository,
	commitStatsRepo repository.ICommitStatsRepository,
	githubGateway gateway.IGithubGateway,
) ISyncCommitsUsecase {
	return &syncCommitsUsecase{
		userRepo:        userRepo,
		rivalRepo:       rivalRepo,
		commitStatsRepo: commitStatsRepo,
		githubGateway:   githubGateway,
	}
}

func (u *syncCommitsUsecase) SyncAllUsers(ctx context.Context) error {
	return u.SyncAllUsersWithDateRange(ctx, nil, nil)
}

func (u *syncCommitsUsecase) SyncAllUsersWithDateRange(ctx context.Context, fromDate, toDate *time.Time) error {
	// 同期が必要なGithubユーザーIDを収集（ユーザー + ライバル）
	syncTargets := make(map[uint64]string) // githubUserID -> username

	log.Println("SyncAllUsers: collecting users to sync...")

	// 全ユーザーを取得
	users, err := u.userRepo.FindAll(ctx)
	if err != nil {
		log.Printf("Failed to get all users: %v", err)
		return err
	}
	for _, user := range users {
		syncTargets[user.GithubUserID] = user.GithubUsername
	}
	log.Printf("Found %d users to sync", len(users))

	// 全ライバルを取得（重複排除済み）
	rivals, err := u.rivalRepo.FindAllDistinctRivals(ctx)
	if err != nil {
		log.Printf("Failed to get all rivals: %v", err)
		return err
	}
	for _, rival := range rivals {
		// ユーザーと重複していない場合のみ追加
		if _, exists := syncTargets[rival.RivalGithubUserID]; !exists {
			syncTargets[rival.RivalGithubUserID] = rival.RivalGithubUsername
		}
	}
	log.Printf("Total %d unique users/rivals to sync", len(syncTargets))

	// 各ユーザーのコミット情報を同期
	var syncErrors []error
	for githubUserID, username := range syncTargets {
		if err := u.SyncUser(ctx, githubUserID, username, fromDate, toDate); err != nil {
			log.Printf("Failed to sync user %s: %v", username, err)
			syncErrors = append(syncErrors, err)
			continue
		}
	}

	if len(syncErrors) > 0 {
		log.Printf("Sync completed with %d errors", len(syncErrors))
	} else {
		log.Println("Sync completed successfully")
	}

	return nil
}

func (u *syncCommitsUsecase) SyncUser(ctx context.Context, githubUserID uint64, githubUsername string, fromDate, toDate *time.Time) error {
	log.Printf("Syncing commits for user: %s", githubUsername)

	// 日付範囲を設定（デフォルトは過去1年）
	now := time.Now()
	from := now.AddDate(-1, 0, 0)
	to := now

	if fromDate != nil {
		from = *fromDate
	}
	if toDate != nil {
		to = *toDate
	}

	// ユーザーの公開リポジトリ一覧を取得
	repos, err := u.githubGateway.GetUserPublicRepos(ctx, githubUsername)
	if err != nil {
		log.Printf("Failed to get repos for %s: %v", githubUsername, err)
		return err
	}

	log.Printf("Found %d public repos for user: %s", len(repos), githubUsername)

	// リポジトリ別・日別のコミット数を集計
	// map[date][repo] = count
	commitsByDateAndRepo := make(map[string]map[string]int)

	for _, repo := range repos {
		commits, err := u.githubGateway.GetRepositoryCommits(ctx, repo.Owner.Login, repo.Name, githubUsername, from, to)
		if err != nil {
			log.Printf("Failed to get commits for %s/%s: %v", repo.Owner.Login, repo.Name, err)
			continue
		}

		if len(commits) == 0 {
			continue
		}

		repoFullName := repo.FullName
		for _, commit := range commits {
			dateStr := commit.Commit.Author.Date.Format("2006-01-02")
			if commitsByDateAndRepo[dateStr] == nil {
				commitsByDateAndRepo[dateStr] = make(map[string]int)
			}
			commitsByDateAndRepo[dateStr][repoFullName]++
		}
	}

	// コミット統計を保存
	var statsList []models.CommitStats
	for dateStr, repos := range commitsByDateAndRepo {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		for repo, count := range repos {
			statsList = append(statsList, models.CommitStats{
				GithubUserID:   githubUserID,
				GithubUsername: githubUsername,
				Date:           date,
				Repository:     repo,
				CommitCount:    count,
			})
		}
	}

	if len(statsList) > 0 {
		for _, s := range statsList {
			log.Printf("  -> %s %s: %d commits", s.Date.Format("2006-01-02"), s.Repository, s.CommitCount)
		}
		if err := u.commitStatsRepo.UpsertBatch(ctx, statsList); err != nil {
			return err
		}
		log.Printf("Saved %d commit stats for user: %s", len(statsList), githubUsername)
	} else {
		log.Printf("No commit stats to save for user: %s", githubUsername)
	}

	return nil
}
