package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
)

// ISignalUsecase シグナルユースケースのインターフェース
type ISignalUsecase interface {
	GetSignals(ctx context.Context, userID uint64, circleID uint64) ([]Signal, error)
	GetRecentSignals(ctx context.Context, userID uint64) ([]Signal, error)
}

// Signal シグナル情報
type Signal struct {
	Type       string
	Date       string
	Usernames  []string
	AvatarURLs []string
	Detail     string
	CircleID   uint64
	CircleName string
}

type signalUsecase struct {
	circleRepo      repository.ICircleRepository
	commitStatsRepo repository.ICommitStatsRepository
}

// NewSignalUsecase コンストラクタ
func NewSignalUsecase(circleRepo repository.ICircleRepository, commitStatsRepo repository.ICommitStatsRepository) ISignalUsecase {
	return &signalUsecase{
		circleRepo:      circleRepo,
		commitStatsRepo: commitStatsRepo,
	}
}

func (u *signalUsecase) GetSignals(ctx context.Context, userID uint64, circleID uint64) ([]Signal, error) {
	circle, err := u.circleRepo.FindByID(ctx, circleID)
	if err != nil {
		return nil, fmt.Errorf("サークルが見つかりません")
	}

	// 権限チェック
	isMember := false
	for _, member := range circle.Members {
		if member.UserID == userID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, fmt.Errorf("このサークルのメンバーではありません")
	}

	return u.detectSignals(ctx, userID, circle)
}

func (u *signalUsecase) GetRecentSignals(ctx context.Context, userID uint64) ([]Signal, error) {
	circles, err := u.circleRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var allSignals []Signal
	for i := range circles {
		signals, err := u.detectSignals(ctx, userID, &circles[i])
		if err != nil {
			continue
		}
		for j := range signals {
			signals[j].CircleID = circles[i].ID
			signals[j].CircleName = circles[i].Name
		}
		allSignals = append(allSignals, signals...)
	}

	// 日付降順でソート
	sort.Slice(allSignals, func(i, j int) bool {
		return allSignals[i].Date > allSignals[j].Date
	})

	// 最大10件
	if len(allSignals) > 10 {
		allSignals = allSignals[:10]
	}

	return allSignals, nil
}

func (u *signalUsecase) detectSignals(ctx context.Context, userID uint64, circle *models.Circle) ([]Signal, error) {
	// メンバーの GithubUserID を収集
	var myGithubUserID uint64
	memberMap := make(map[uint64]models.User) // githubUserID -> User
	var githubUserIDs []uint64

	for _, member := range circle.Members {
		githubUserIDs = append(githubUserIDs, member.User.GithubUserID)
		memberMap[member.User.GithubUserID] = member.User
		if member.UserID == userID {
			myGithubUserID = member.User.GithubUserID
		}
	}

	if myGithubUserID == 0 {
		return nil, fmt.Errorf("このサークルのメンバーではありません")
	}

	// 直近7日間のデータ取得
	now := time.Now()
	startDate := now.AddDate(0, 0, -7)
	stats, err := u.commitStatsRepo.FindByGithubUserIDsAndDateRange(ctx, githubUserIDs, startDate, now)
	if err != nil {
		return nil, err
	}

	// ユーザーごと・日付ごとにデータを整理
	type dayData struct {
		hasCommit bool
		hours     map[int]bool
		languages map[string]bool
	}
	// userGithubID -> date -> dayData
	userDayMap := make(map[uint64]map[string]*dayData)

	for _, s := range stats {
		if userDayMap[s.GithubUserID] == nil {
			userDayMap[s.GithubUserID] = make(map[string]*dayData)
		}
		dateStr := s.Date.Format("2006-01-02")
		dd := userDayMap[s.GithubUserID][dateStr]
		if dd == nil {
			dd = &dayData{
				hours:     make(map[int]bool),
				languages: make(map[string]bool),
			}
			userDayMap[s.GithubUserID][dateStr] = dd
		}
		dd.hasCommit = true
		if s.PrimaryHour != nil {
			dd.hours[*s.PrimaryHour] = true
		}
		if s.Language != "" {
			dd.languages[s.Language] = true
		}
	}

	myDays := userDayMap[myGithubUserID]
	if myDays == nil {
		return []Signal{}, nil
	}

	var signals []Signal
	// 重複排除用
	type signalKey struct {
		typ    string
		date   string
		detail string
	}
	seen := make(map[signalKey]bool)

	for dateStr, myDay := range myDays {
		if !myDay.hasCommit {
			continue
		}

		for otherGithubUserID, otherDays := range userDayMap {
			if otherGithubUserID == myGithubUserID {
				continue
			}

			otherDay, exists := otherDays[dateStr]
			if !exists || !otherDay.hasCommit {
				continue
			}

			otherUser := memberMap[otherGithubUserID]

			// 同日コミット
			key := signalKey{typ: "same_day", date: dateStr}
			if !seen[key] {
				var usernames []string
				var avatarURLs []string
				for uid, days := range userDayMap {
					if uid == myGithubUserID {
						continue
					}
					if d, ok := days[dateStr]; ok && d.hasCommit {
						u := memberMap[uid]
						usernames = append(usernames, u.GithubUsername)
						avatarURLs = append(avatarURLs, u.AvatarURL)
					}
				}
				signals = append(signals, Signal{
					Type:       "same_day",
					Date:       dateStr,
					Usernames:  usernames,
					AvatarURLs: avatarURLs,
					Detail:     "同じ日にコミット",
				})
				seen[key] = true
			}

			// 同時間帯コミット（±1時間）
			for myHour := range myDay.hours {
				for otherHour := range otherDay.hours {
					diff := myHour - otherHour
					if diff < 0 {
						diff = -diff
					}
					if diff <= 1 {
						detail := fmt.Sprintf("%d時台", otherHour)
						key := signalKey{typ: "same_hour", date: dateStr, detail: detail}
						if !seen[key] {
							signals = append(signals, Signal{
								Type:       "same_hour",
								Date:       dateStr,
								Usernames:  []string{otherUser.GithubUsername},
								AvatarURLs: []string{otherUser.AvatarURL},
								Detail:     detail,
							})
							seen[key] = true
						}
					}
				}
			}

			// 同言語使用
			for lang := range myDay.languages {
				if otherDay.languages[lang] {
					key := signalKey{typ: "same_language", date: dateStr, detail: lang}
					if !seen[key] {
						signals = append(signals, Signal{
							Type:       "same_language",
							Date:       dateStr,
							Usernames:  []string{otherUser.GithubUsername},
							AvatarURLs: []string{otherUser.AvatarURL},
							Detail:     lang,
						})
						seen[key] = true
					}
				}
			}
		}
	}

	// 日付降順でソート
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].Date > signals[j].Date
	})

	return signals, nil
}
