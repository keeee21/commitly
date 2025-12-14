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

// userMockUserRepository テスト用のモックリポジトリ
type userMockUserRepository struct {
	FindByGithubUserIDFunc func(ctx context.Context, githubUserID uint64) (*models.User, error)
	CreateFunc             func(ctx context.Context, user *models.User) error
	UpdateFunc             func(ctx context.Context, user *models.User) error
}

func (m *userMockUserRepository) FindByGithubUserID(ctx context.Context, githubUserID uint64) (*models.User, error) {
	if m.FindByGithubUserIDFunc != nil {
		return m.FindByGithubUserIDFunc(ctx, githubUserID)
	}
	return nil, nil
}

func (m *userMockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *userMockUserRepository) Update(ctx context.Context, user *models.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *userMockUserRepository) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	return nil, nil
}

func (m *userMockUserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	return nil, nil
}

// userMockGithubGateway テスト用のモックゲートウェイ
type userMockGithubGateway struct{}

func (m *userMockGithubGateway) GetUser(ctx context.Context, username string) (*gateway.GithubUser, error) {
	return nil, nil
}

func (m *userMockGithubGateway) GetUserEvents(ctx context.Context, username string, page int) ([]gateway.GithubEvent, error) {
	return nil, nil
}

func (m *userMockGithubGateway) GetUserPublicRepos(ctx context.Context, username string) ([]gateway.GithubRepo, error) {
	return nil, nil
}

func (m *userMockGithubGateway) GetUserContributions(ctx context.Context, username string, from, to string) ([]gateway.ContributionDay, error) {
	return nil, nil
}

func (m *userMockGithubGateway) GetRepositoryCommits(ctx context.Context, owner, repo, author string, since, until time.Time) ([]gateway.RepositoryCommit, error) {
	return nil, nil
}

func TestGetOrCreateUser_NewUser(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := &userMockUserRepository{
		FindByGithubUserIDFunc: func(ctx context.Context, githubUserID uint64) (*models.User, error) {
			return nil, errors.New("not found")
		},
		CreateFunc: func(ctx context.Context, user *models.User) error {
			user.ID = 1
			return nil
		},
	}
	mockGithubGateway := &userMockGithubGateway{}

	usecase := NewUserUsecase(mockUserRepo, mockGithubGateway)

	user, err := usecase.GetOrCreateUser(ctx, 12345, "testuser", "test@example.com", "https://avatar.url")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint64(1), user.ID)
	assert.Equal(t, uint64(12345), user.GithubUserID)
	assert.Equal(t, "testuser", user.GithubUsername)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "https://avatar.url", user.AvatarURL)
}

func TestGetOrCreateUser_ExistingUser(t *testing.T) {
	ctx := context.Background()
	existingUser := &models.User{
		ID:             1,
		GithubUserID:   12345,
		GithubUsername: "oldusername",
		Email:          "old@example.com",
		AvatarURL:      "https://old.avatar.url",
	}

	mockUserRepo := &userMockUserRepository{
		FindByGithubUserIDFunc: func(ctx context.Context, githubUserID uint64) (*models.User, error) {
			return existingUser, nil
		},
		UpdateFunc: func(ctx context.Context, user *models.User) error {
			return nil
		},
	}
	mockGithubGateway := &userMockGithubGateway{}

	usecase := NewUserUsecase(mockUserRepo, mockGithubGateway)

	user, err := usecase.GetOrCreateUser(ctx, 12345, "newusername", "new@example.com", "https://new.avatar.url")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint64(1), user.ID)
	assert.Equal(t, "newusername", user.GithubUsername)
	assert.Equal(t, "new@example.com", user.Email)
	assert.Equal(t, "https://new.avatar.url", user.AvatarURL)
}

func TestGetOrCreateUser_CreateError(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := &userMockUserRepository{
		FindByGithubUserIDFunc: func(ctx context.Context, githubUserID uint64) (*models.User, error) {
			return nil, errors.New("not found")
		},
		CreateFunc: func(ctx context.Context, user *models.User) error {
			return errors.New("database error")
		},
	}
	mockGithubGateway := &userMockGithubGateway{}

	usecase := NewUserUsecase(mockUserRepo, mockGithubGateway)

	user, err := usecase.GetOrCreateUser(ctx, 12345, "testuser", "test@example.com", "https://avatar.url")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "database error", err.Error())
}

func TestGetUserByGithubUserID_Success(t *testing.T) {
	ctx := context.Background()
	expectedUser := &models.User{
		ID:             1,
		GithubUserID:   12345,
		GithubUsername: "testuser",
	}

	mockUserRepo := &userMockUserRepository{
		FindByGithubUserIDFunc: func(ctx context.Context, githubUserID uint64) (*models.User, error) {
			return expectedUser, nil
		},
	}
	mockGithubGateway := &userMockGithubGateway{}

	usecase := NewUserUsecase(mockUserRepo, mockGithubGateway)

	user, err := usecase.GetUserByGithubUserID(ctx, 12345)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserByGithubUserID_NotFound(t *testing.T) {
	ctx := context.Background()
	mockUserRepo := &userMockUserRepository{
		FindByGithubUserIDFunc: func(ctx context.Context, githubUserID uint64) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}
	mockGithubGateway := &userMockGithubGateway{}

	usecase := NewUserUsecase(mockUserRepo, mockGithubGateway)

	user, err := usecase.GetUserByGithubUserID(ctx, 99999)

	assert.Error(t, err)
	assert.Nil(t, user)
}
