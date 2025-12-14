package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/keeee21/commitly/api/gateway"
	"github.com/keeee21/commitly/api/models"
	"github.com/stretchr/testify/assert"
)

// syncMockUserRepository テスト用のモックリポジトリ
type syncMockUserRepository struct {
	FindAllFunc func(ctx context.Context) ([]models.User, error)
}

func (m *syncMockUserRepository) FindByGithubUserID(ctx context.Context, githubUserID uint64) (*models.User, error) {
	return nil, nil
}

func (m *syncMockUserRepository) Create(ctx context.Context, user *models.User) error {
	return nil
}

func (m *syncMockUserRepository) Update(ctx context.Context, user *models.User) error {
	return nil
}

func (m *syncMockUserRepository) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	return nil, nil
}

func (m *syncMockUserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

// syncMockRivalRepository テスト用のモックリポジトリ
type syncMockRivalRepository struct {
	FindAllDistinctRivalsFunc func(ctx context.Context) ([]models.Rival, error)
}

func (m *syncMockRivalRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.Rival, error) {
	return nil, nil
}

func (m *syncMockRivalRepository) Create(ctx context.Context, rival *models.Rival) error {
	return nil
}

func (m *syncMockRivalRepository) FindByID(ctx context.Context, id uint64) (*models.Rival, error) {
	return nil, nil
}

func (m *syncMockRivalRepository) Delete(ctx context.Context, id uint64) error {
	return nil
}

func (m *syncMockRivalRepository) CountByUserID(ctx context.Context, userID uint64) (int64, error) {
	return 0, nil
}

func (m *syncMockRivalRepository) ExistsByUserIDAndRivalGithubUserID(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error) {
	return false, nil
}

func (m *syncMockRivalRepository) FindAllDistinctRivals(ctx context.Context) ([]models.Rival, error) {
	if m.FindAllDistinctRivalsFunc != nil {
		return m.FindAllDistinctRivalsFunc(ctx)
	}
	return nil, nil
}

// syncMockCommitStatsRepository テスト用のモックリポジトリ
type syncMockCommitStatsRepository struct {
	UpsertBatchFunc func(ctx context.Context, statsList []models.CommitStats) error
}

func (m *syncMockCommitStatsRepository) FindByGithubUserIDAndDateRange(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	return nil, nil
}

func (m *syncMockCommitStatsRepository) FindByGithubUserIDsAndDateRange(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	return nil, nil
}

func (m *syncMockCommitStatsRepository) Upsert(ctx context.Context, stats *models.CommitStats) error {
	return nil
}

func (m *syncMockCommitStatsRepository) UpsertBatch(ctx context.Context, statsList []models.CommitStats) error {
	if m.UpsertBatchFunc != nil {
		return m.UpsertBatchFunc(ctx, statsList)
	}
	return nil
}

// syncMockGithubGateway テスト用のモックゲートウェイ
type syncMockGithubGateway struct {
	GetUserPublicReposFunc   func(ctx context.Context, username string) ([]gateway.GithubRepo, error)
	GetRepositoryCommitsFunc func(ctx context.Context, owner, repo, author string, since, until time.Time) ([]gateway.RepositoryCommit, error)
}

func (m *syncMockGithubGateway) GetUser(ctx context.Context, username string) (*gateway.GithubUser, error) {
	return nil, nil
}

func (m *syncMockGithubGateway) GetUserEvents(ctx context.Context, username string, page int) ([]gateway.GithubEvent, error) {
	return nil, nil
}

func (m *syncMockGithubGateway) GetUserPublicRepos(ctx context.Context, username string) ([]gateway.GithubRepo, error) {
	if m.GetUserPublicReposFunc != nil {
		return m.GetUserPublicReposFunc(ctx, username)
	}
	return nil, nil
}

func (m *syncMockGithubGateway) GetUserContributions(ctx context.Context, username string, from, to string) ([]gateway.ContributionDay, error) {
	return nil, nil
}

func (m *syncMockGithubGateway) GetRepositoryCommits(ctx context.Context, owner, repo, author string, since, until time.Time) ([]gateway.RepositoryCommit, error) {
	if m.GetRepositoryCommitsFunc != nil {
		return m.GetRepositoryCommitsFunc(ctx, owner, repo, author, since, until)
	}
	return nil, nil
}

func TestSyncAllUsers_Success(t *testing.T) {
	ctx := context.Background()

	users := []models.User{
		{ID: 1, GithubUserID: 100, GithubUsername: "user1"},
	}

	rivals := []models.Rival{
		{RivalGithubUserID: 200, RivalGithubUsername: "rival1"},
	}

	repos := []gateway.GithubRepo{
		{Name: "repo1", FullName: "user1/repo1", Owner: struct {
			Login string `json:"login"`
		}{Login: "user1"}},
	}

	now := time.Now()
	commits := []gateway.RepositoryCommit{
		{SHA: "abc123", Commit: struct {
			Author struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"author"`
			Message string `json:"message"`
		}{Author: struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		}{Date: now}}},
	}

	mockUserRepo := &syncMockUserRepository{
		FindAllFunc: func(ctx context.Context) ([]models.User, error) {
			return users, nil
		},
	}

	mockRivalRepo := &syncMockRivalRepository{
		FindAllDistinctRivalsFunc: func(ctx context.Context) ([]models.Rival, error) {
			return rivals, nil
		},
	}

	mockCommitStatsRepo := &syncMockCommitStatsRepository{
		UpsertBatchFunc: func(ctx context.Context, statsList []models.CommitStats) error {
			return nil
		},
	}

	mockGithubGateway := &syncMockGithubGateway{
		GetUserPublicReposFunc: func(ctx context.Context, username string) ([]gateway.GithubRepo, error) {
			return repos, nil
		},
		GetRepositoryCommitsFunc: func(ctx context.Context, owner, repo, author string, since, until time.Time) ([]gateway.RepositoryCommit, error) {
			return commits, nil
		},
	}

	usecase := NewSyncCommitsUsecase(mockUserRepo, mockRivalRepo, mockCommitStatsRepo, mockGithubGateway)
	err := usecase.SyncAllUsers(ctx)

	assert.NoError(t, err)
}

func TestSyncAllUsers_NoUsers(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := &syncMockUserRepository{
		FindAllFunc: func(ctx context.Context) ([]models.User, error) {
			return []models.User{}, nil
		},
	}

	mockRivalRepo := &syncMockRivalRepository{
		FindAllDistinctRivalsFunc: func(ctx context.Context) ([]models.Rival, error) {
			return []models.Rival{}, nil
		},
	}

	mockCommitStatsRepo := &syncMockCommitStatsRepository{}
	mockGithubGateway := &syncMockGithubGateway{}

	usecase := NewSyncCommitsUsecase(mockUserRepo, mockRivalRepo, mockCommitStatsRepo, mockGithubGateway)
	err := usecase.SyncAllUsers(ctx)

	assert.NoError(t, err)
}

func TestSyncAllUsers_UserRepoError(t *testing.T) {
	ctx := context.Background()

	mockUserRepo := &syncMockUserRepository{
		FindAllFunc: func(ctx context.Context) ([]models.User, error) {
			return nil, errors.New("database error")
		},
	}

	mockRivalRepo := &syncMockRivalRepository{}
	mockCommitStatsRepo := &syncMockCommitStatsRepository{}
	mockGithubGateway := &syncMockGithubGateway{}

	usecase := NewSyncCommitsUsecase(mockUserRepo, mockRivalRepo, mockCommitStatsRepo, mockGithubGateway)
	err := usecase.SyncAllUsers(ctx)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
}

func TestSyncUser_Success(t *testing.T) {
	ctx := context.Background()

	repos := []gateway.GithubRepo{
		{Name: "repo1", FullName: "testuser/repo1", Owner: struct {
			Login string `json:"login"`
		}{Login: "testuser"}},
		{Name: "repo2", FullName: "testuser/repo2", Owner: struct {
			Login string `json:"login"`
		}{Login: "testuser"}},
	}

	now := time.Now()
	commits := []gateway.RepositoryCommit{
		{SHA: "abc123", Commit: struct {
			Author struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"author"`
			Message string `json:"message"`
		}{Author: struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		}{Date: now}}},
		{SHA: "def456", Commit: struct {
			Author struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"author"`
			Message string `json:"message"`
		}{Author: struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		}{Date: now.AddDate(0, 0, -1)}}},
	}

	var savedStats []models.CommitStats
	mockCommitStatsRepo := &syncMockCommitStatsRepository{
		UpsertBatchFunc: func(ctx context.Context, statsList []models.CommitStats) error {
			savedStats = statsList
			return nil
		},
	}

	mockGithubGateway := &syncMockGithubGateway{
		GetUserPublicReposFunc: func(ctx context.Context, username string) ([]gateway.GithubRepo, error) {
			return repos, nil
		},
		GetRepositoryCommitsFunc: func(ctx context.Context, owner, repo, author string, since, until time.Time) ([]gateway.RepositoryCommit, error) {
			return commits, nil
		},
	}

	mockUserRepo := &syncMockUserRepository{}
	mockRivalRepo := &syncMockRivalRepository{}

	usecase := NewSyncCommitsUsecase(mockUserRepo, mockRivalRepo, mockCommitStatsRepo, mockGithubGateway)
	err := usecase.SyncUser(ctx, 100, "testuser", nil, nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, savedStats)
}

func TestSyncUser_NoRepos(t *testing.T) {
	ctx := context.Background()

	mockGithubGateway := &syncMockGithubGateway{
		GetUserPublicReposFunc: func(ctx context.Context, username string) ([]gateway.GithubRepo, error) {
			return []gateway.GithubRepo{}, nil
		},
	}

	mockUserRepo := &syncMockUserRepository{}
	mockRivalRepo := &syncMockRivalRepository{}
	mockCommitStatsRepo := &syncMockCommitStatsRepository{}

	usecase := NewSyncCommitsUsecase(mockUserRepo, mockRivalRepo, mockCommitStatsRepo, mockGithubGateway)
	err := usecase.SyncUser(ctx, 100, "testuser", nil, nil)

	assert.NoError(t, err)
}

func TestSyncUser_WithDateRange(t *testing.T) {
	ctx := context.Background()

	repos := []gateway.GithubRepo{
		{Name: "repo1", FullName: "testuser/repo1", Owner: struct {
			Login string `json:"login"`
		}{Login: "testuser"}},
	}

	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	var capturedSince, capturedUntil time.Time
	mockGithubGateway := &syncMockGithubGateway{
		GetUserPublicReposFunc: func(ctx context.Context, username string) ([]gateway.GithubRepo, error) {
			return repos, nil
		},
		GetRepositoryCommitsFunc: func(ctx context.Context, owner, repo, author string, since, until time.Time) ([]gateway.RepositoryCommit, error) {
			capturedSince = since
			capturedUntil = until
			return []gateway.RepositoryCommit{}, nil
		},
	}

	mockUserRepo := &syncMockUserRepository{}
	mockRivalRepo := &syncMockRivalRepository{}
	mockCommitStatsRepo := &syncMockCommitStatsRepository{}

	usecase := NewSyncCommitsUsecase(mockUserRepo, mockRivalRepo, mockCommitStatsRepo, mockGithubGateway)
	err := usecase.SyncUser(ctx, 100, "testuser", &from, &to)

	assert.NoError(t, err)
	assert.Equal(t, from, capturedSince)
	assert.Equal(t, to, capturedUntil)
}

func TestSyncUser_GitHubAPIError(t *testing.T) {
	ctx := context.Background()

	mockGithubGateway := &syncMockGithubGateway{
		GetUserPublicReposFunc: func(ctx context.Context, username string) ([]gateway.GithubRepo, error) {
			return nil, errors.New("GitHub API error")
		},
	}

	mockUserRepo := &syncMockUserRepository{}
	mockRivalRepo := &syncMockRivalRepository{}
	mockCommitStatsRepo := &syncMockCommitStatsRepository{}

	usecase := NewSyncCommitsUsecase(mockUserRepo, mockRivalRepo, mockCommitStatsRepo, mockGithubGateway)
	err := usecase.SyncUser(ctx, 100, "testuser", nil, nil)

	assert.Error(t, err)
	assert.Equal(t, "GitHub API error", err.Error())
}
