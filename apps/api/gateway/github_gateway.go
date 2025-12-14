package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// IGithubGateway GitHub APIゲートウェイのインターフェース
type IGithubGateway interface {
	GetUser(ctx context.Context, username string) (*GithubUser, error)
	GetUserEvents(ctx context.Context, username string, page int) ([]GithubEvent, error)
	GetUserPublicRepos(ctx context.Context, username string) ([]GithubRepo, error)
	GetUserContributions(ctx context.Context, username string, from, to string) ([]ContributionDay, error)
	GetRepositoryCommits(ctx context.Context, owner, repo, author string, since, until time.Time) ([]RepositoryCommit, error)
}

// GithubUser GitHubユーザー情報
type GithubUser struct {
	ID        uint64 `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

// GithubRepo GitHubリポジトリ情報
type GithubRepo struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
	Private bool `json:"private"`
}

// GithubEvent GitHubイベント情報
type GithubEvent struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Repo      struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	} `json:"repo"`
	Payload json.RawMessage `json:"payload"`
}

// PushEventPayload PushEventのペイロード
type PushEventPayload struct {
	Size    int `json:"size"`
	Commits []struct {
		SHA string `json:"sha"`
	} `json:"commits"`
}

// ContributionDay 日別コントリビューション
type ContributionDay struct {
	Date              string `json:"date"`
	ContributionCount int    `json:"contributionCount"`
}

// RepositoryCommit リポジトリのコミット情報
type RepositoryCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
}

type githubGateway struct {
	token      string
	httpClient *http.Client
}

// NewGithubGateway コンストラクタ
func NewGithubGateway(token string) IGithubGateway {
	return &githubGateway{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (g *githubGateway) doRequest(ctx context.Context, url string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if g.token != "" {
		req.Header.Set("Authorization", "Bearer "+g.token)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found: %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Github API error: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (g *githubGateway) GetUser(ctx context.Context, username string) (*GithubUser, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	var user GithubUser
	if err := g.doRequest(ctx, url, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (g *githubGateway) GetUserEvents(ctx context.Context, username string, page int) ([]GithubEvent, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events/public?per_page=100&page=%d", username, page)
	var events []GithubEvent
	if err := g.doRequest(ctx, url, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func (g *githubGateway) GetUserPublicRepos(ctx context.Context, username string) ([]GithubRepo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?type=public&per_page=100", username)
	var repos []GithubRepo
	if err := g.doRequest(ctx, url, &repos); err != nil {
		return nil, err
	}
	return repos, nil
}

// GetUserContributions GraphQL APIでユーザーのコントリビューションを取得
func (g *githubGateway) GetUserContributions(ctx context.Context, username string, from, to string) ([]ContributionDay, error) {
	query := fmt.Sprintf(`{
		user(login: "%s") {
			contributionsCollection(from: "%s", to: "%s") {
				contributionCalendar {
					weeks {
						contributionDays {
							date
							contributionCount
						}
					}
				}
			}
		}
	}`, username, from, to)

	reqBody, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.github.com/graphql", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.token)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GraphQL API error: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			User struct {
				ContributionsCollection struct {
					ContributionCalendar struct {
						Weeks []struct {
							ContributionDays []ContributionDay `json:"contributionDays"`
						} `json:"weeks"`
					} `json:"contributionCalendar"`
				} `json:"contributionsCollection"`
			} `json:"user"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	var days []ContributionDay
	for _, week := range result.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		days = append(days, week.ContributionDays...)
	}

	return days, nil
}

// GetRepositoryCommits リポジトリのコミット履歴を取得
func (g *githubGateway) GetRepositoryCommits(ctx context.Context, owner, repo, author string, since, until time.Time) ([]RepositoryCommit, error) {
	var allCommits []RepositoryCommit

	for page := 1; page <= 10; page++ { // 最大10ページ
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/commits?author=%s&since=%s&until=%s&per_page=100&page=%d",
			owner, repo, author,
			since.Format(time.RFC3339),
			until.Format(time.RFC3339),
			page,
		)

		var commits []RepositoryCommit
		if err := g.doRequest(ctx, url, &commits); err != nil {
			// 404やエラーの場合は空で返す（フォークなどでアクセスできない場合がある）
			if page == 1 {
				return nil, nil
			}
			break
		}

		if len(commits) == 0 {
			break
		}

		allCommits = append(allCommits, commits...)

		if len(commits) < 100 {
			break
		}
	}

	return allCommits, nil
}

// ExtractCommitsFromEvents イベントからコミット情報を抽出する
func ExtractCommitsFromEvents(events []GithubEvent) map[string]map[string]int {
	// result[date][repo] = commitCount
	result := make(map[string]map[string]int)

	for _, event := range events {
		if event.Type != "PushEvent" {
			continue
		}

		dateStr := event.CreatedAt.Format("2006-01-02")
		repoName := event.Repo.Name

		if result[dateStr] == nil {
			result[dateStr] = make(map[string]int)
		}

		var payload PushEventPayload
		if err := json.Unmarshal(event.Payload, &payload); err != nil {
			fmt.Printf("DEBUG: Failed to unmarshal payload: %v\n", err)
			fmt.Printf("DEBUG: Raw payload: %s\n", string(event.Payload))
			continue
		}

		// payload.Sizeが0の場合はcommits配列の長さを使用
		commitCount := payload.Size
		if commitCount == 0 {
			commitCount = len(payload.Commits)
		}

		// デバッグ出力
		if commitCount == 0 {
			fmt.Printf("DEBUG: Zero commits - Size=%d, Commits=%d, Raw=%s\n",
				payload.Size, len(payload.Commits), string(event.Payload)[:min(200, len(event.Payload))])
		}

		result[dateStr][repoName] += commitCount
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
