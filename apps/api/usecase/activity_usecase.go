package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/keeee21/commitly/api/dto"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
)

// IActivityUsecase アクティビティユースケースのインターフェース
type IActivityUsecase interface {
	GetActivityStream(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.ActivityStreamResponse, error)
	GetRhythm(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.RhythmResponse, error)
}

type activityUsecase struct {
	commitStatsRepo repository.ICommitStatsRepository
}

// NewActivityUsecase コンストラクタ
func NewActivityUsecase(commitStatsRepo repository.ICommitStatsRepository) IActivityUsecase {
	return &activityUsecase{
		commitStatsRepo: commitStatsRepo,
	}
}

func (u *activityUsecase) GetActivityStream(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.ActivityStreamResponse, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -6)

	githubUserIDs, userInfoMap := u.collectUserInfo(user, rivals)

	stats, err := u.commitStatsRepo.FindByGithubUserIDsAndDateRange(ctx, githubUserIDs, startDate, now)
	if err != nil {
		return nil, fmt.Errorf("アクティビティデータの取得に失敗しました")
	}

	activities := make([]dto.ActivityItem, 0, len(stats))
	for _, stat := range stats {
		info := userInfoMap[stat.GithubUserID]
		activities = append(activities, dto.ActivityItem{
			GithubUsername: info.username,
			AvatarURL:      info.avatarURL,
			Repository:     stat.Repository,
			CommitCount:    stat.CommitCount,
			Date:           stat.Date.Format("2006-01-02"),
		})
	}

	sort.Slice(activities, func(i, j int) bool {
		return activities[i].Date > activities[j].Date
	})

	return &dto.ActivityStreamResponse{
		Activities: activities,
	}, nil
}

func (u *activityUsecase) GetRhythm(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.RhythmResponse, error) {
	now := time.Now()
	startDate := now.AddDate(0, 0, -6)

	githubUserIDs, userInfoMap := u.collectUserInfo(user, rivals)

	stats, err := u.commitStatsRepo.FindByGithubUserIDsAndDateRange(ctx, githubUserIDs, startDate, now)
	if err != nil {
		return nil, fmt.Errorf("リズムデータの取得に失敗しました")
	}

	// ユーザーごとの曜日別コミット有無を集計
	type weekdaySet struct {
		days [7]bool // 0=Sun, 1=Mon, ..., 6=Sat
	}
	userWeekdays := make(map[uint64]*weekdaySet)
	for _, id := range githubUserIDs {
		userWeekdays[id] = &weekdaySet{}
	}

	for _, stat := range stats {
		if ws, ok := userWeekdays[stat.GithubUserID]; ok {
			weekday := stat.Date.Weekday() // 0=Sunday
			ws.days[weekday] = true
		}
	}

	users := make([]dto.UserRhythm, 0, len(githubUserIDs))

	// 自分を先頭にするため、まず自分、次にライバルの順で処理
	orderedIDs := []uint64{user.GithubUserID}
	for _, rival := range rivals {
		orderedIDs = append(orderedIDs, rival.RivalGithubUserID)
	}

	for _, id := range orderedIDs {
		ws := userWeekdays[id]
		info := userInfoMap[id]

		rhythm := dto.WeeklyRhythm{
			Mon: ws.days[time.Monday],
			Tue: ws.days[time.Tuesday],
			Wed: ws.days[time.Wednesday],
			Thu: ws.days[time.Thursday],
			Fri: ws.days[time.Friday],
			Sat: ws.days[time.Saturday],
			Sun: ws.days[time.Sunday],
		}

		activeDays := 0
		for _, d := range ws.days {
			if d {
				activeDays++
			}
		}

		weekdayCount := 0
		for _, d := range []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday} {
			if ws.days[d] {
				weekdayCount++
			}
		}
		weekendCount := 0
		if ws.days[time.Saturday] {
			weekendCount++
		}
		if ws.days[time.Sunday] {
			weekendCount++
		}

		patternLabel := classifyPattern(activeDays, weekdayCount, weekendCount)

		users = append(users, dto.UserRhythm{
			GithubUsername: info.username,
			AvatarURL:      info.avatarURL,
			PatternLabel:   patternLabel,
			WeeklyRhythm:   rhythm,
		})
	}

	return &dto.RhythmResponse{
		Users:  users,
		Period: fmt.Sprintf("%s/%s", startDate.Format("2006-01-02"), now.Format("2006-01-02")),
	}, nil
}

type userInfo struct {
	username  string
	avatarURL string
}

func (u *activityUsecase) collectUserInfo(user *models.User, rivals []models.Rival) ([]uint64, map[uint64]userInfo) {
	githubUserIDs := []uint64{user.GithubUserID}
	userInfoMap := map[uint64]userInfo{
		user.GithubUserID: {
			username:  user.GithubUsername,
			avatarURL: user.AvatarURL,
		},
	}

	for _, rival := range rivals {
		githubUserIDs = append(githubUserIDs, rival.RivalGithubUserID)
		userInfoMap[rival.RivalGithubUserID] = userInfo{
			username:  rival.RivalGithubUsername,
			avatarURL: rival.RivalAvatarURL,
		}
	}

	return githubUserIDs, userInfoMap
}

func classifyPattern(activeDays, weekdayCount, weekendCount int) string {
	if activeDays >= 5 {
		return "安定型"
	}
	if weekdayCount <= 2 && weekendCount >= 1 {
		return "週末型"
	}
	return "バースト型"
}
