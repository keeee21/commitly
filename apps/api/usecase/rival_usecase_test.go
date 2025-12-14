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

// rivalMockRivalRepository テスト用のモックリポジトリ
type rivalMockRivalRepository struct {
	FindByUserIDFunc                       func(ctx context.Context, userID uint64) ([]models.Rival, error)
	FindByIDFunc                           func(ctx context.Context, id uint64) (*models.Rival, error)
	CreateFunc                             func(ctx context.Context, rival *models.Rival) error
	DeleteFunc                             func(ctx context.Context, id uint64) error
	CountByUserIDFunc                      func(ctx context.Context, userID uint64) (int64, error)
	ExistsByUserIDAndRivalGithubUserIDFunc func(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error)
}

func (m *rivalMockRivalRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.Rival, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *rivalMockRivalRepository) FindByID(ctx context.Context, id uint64) (*models.Rival, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *rivalMockRivalRepository) Create(ctx context.Context, rival *models.Rival) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, rival)
	}
	return nil
}

func (m *rivalMockRivalRepository) Delete(ctx context.Context, id uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *rivalMockRivalRepository) CountByUserID(ctx context.Context, userID uint64) (int64, error) {
	if m.CountByUserIDFunc != nil {
		return m.CountByUserIDFunc(ctx, userID)
	}
	return 0, nil
}

func (m *rivalMockRivalRepository) ExistsByUserIDAndRivalGithubUserID(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error) {
	if m.ExistsByUserIDAndRivalGithubUserIDFunc != nil {
		return m.ExistsByUserIDAndRivalGithubUserIDFunc(ctx, userID, rivalGithubUserID)
	}
	return false, nil
}

func (m *rivalMockRivalRepository) FindAllDistinctRivals(ctx context.Context) ([]models.Rival, error) {
	return nil, nil
}

// rivalMockGithubGateway テスト用のモックゲートウェイ
type rivalMockGithubGateway struct {
	GetUserFunc func(ctx context.Context, username string) (*gateway.GithubUser, error)
}

func (m *rivalMockGithubGateway) GetUser(ctx context.Context, username string) (*gateway.GithubUser, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, username)
	}
	return nil, nil
}

func (m *rivalMockGithubGateway) GetUserEvents(ctx context.Context, username string, page int) ([]gateway.GithubEvent, error) {
	return nil, nil
}

func (m *rivalMockGithubGateway) GetUserPublicRepos(ctx context.Context, username string) ([]gateway.GithubRepo, error) {
	return nil, nil
}

func (m *rivalMockGithubGateway) GetUserContributions(ctx context.Context, username string, from, to string) ([]gateway.ContributionDay, error) {
	return nil, nil
}

func (m *rivalMockGithubGateway) GetRepositoryCommits(ctx context.Context, owner, repo, author string, since, until time.Time) ([]gateway.RepositoryCommit, error) {
	return nil, nil
}

func TestGetRivals_Success(t *testing.T) {
	ctx := context.Background()
	expectedRivals := []models.Rival{
		{ID: 1, UserID: 1, RivalGithubUserID: 100, RivalGithubUsername: "rival1"},
		{ID: 2, UserID: 1, RivalGithubUserID: 200, RivalGithubUsername: "rival2"},
	}

	mockRivalRepo := &rivalMockRivalRepository{
		FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return expectedRivals, nil
		},
	}
	mockGithubGateway := &rivalMockGithubGateway{}

	usecase := NewRivalUsecase(mockRivalRepo, mockGithubGateway)

	rivals, err := usecase.GetRivals(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, rivals, 2)
	assert.Equal(t, expectedRivals, rivals)
}

func TestGetRivals_Empty(t *testing.T) {
	ctx := context.Background()
	mockRivalRepo := &rivalMockRivalRepository{
		FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return []models.Rival{}, nil
		},
	}
	mockGithubGateway := &rivalMockGithubGateway{}

	usecase := NewRivalUsecase(mockRivalRepo, mockGithubGateway)

	rivals, err := usecase.GetRivals(ctx, 1)

	assert.NoError(t, err)
	assert.Empty(t, rivals)
}

func TestAddRival_Success(t *testing.T) {
	ctx := context.Background()
	mockRivalRepo := &rivalMockRivalRepository{
		CountByUserIDFunc: func(ctx context.Context, userID uint64) (int64, error) {
			return 0, nil
		},
		ExistsByUserIDAndRivalGithubUserIDFunc: func(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error) {
			return false, nil
		},
		CreateFunc: func(ctx context.Context, rival *models.Rival) error {
			rival.ID = 1
			return nil
		},
	}
	mockGithubGateway := &rivalMockGithubGateway{
		GetUserFunc: func(ctx context.Context, username string) (*gateway.GithubUser, error) {
			return &gateway.GithubUser{
				ID:        12345,
				Login:     "rivaluser",
				AvatarURL: "https://avatar.url",
			}, nil
		},
	}

	usecase := NewRivalUsecase(mockRivalRepo, mockGithubGateway)

	rival, err := usecase.AddRival(ctx, 1, "rivaluser")

	assert.NoError(t, err)
	assert.NotNil(t, rival)
	assert.Equal(t, uint64(1), rival.ID)
	assert.Equal(t, uint64(12345), rival.RivalGithubUserID)
	assert.Equal(t, "rivaluser", rival.RivalGithubUsername)
}

func TestAddRival_MaxRivalsReached(t *testing.T) {
	ctx := context.Background()
	mockRivalRepo := &rivalMockRivalRepository{
		CountByUserIDFunc: func(ctx context.Context, userID uint64) (int64, error) {
			return int64(models.MaxRivalsForFreePlan), nil
		},
	}
	mockGithubGateway := &rivalMockGithubGateway{}

	usecase := NewRivalUsecase(mockRivalRepo, mockGithubGateway)

	rival, err := usecase.AddRival(ctx, 1, "rivaluser")

	assert.Error(t, err)
	assert.Nil(t, rival)
	assert.Contains(t, err.Error(), "上限")
}

func TestAddRival_GithubUserNotFound(t *testing.T) {
	ctx := context.Background()
	mockRivalRepo := &rivalMockRivalRepository{
		CountByUserIDFunc: func(ctx context.Context, userID uint64) (int64, error) {
			return 0, nil
		},
	}
	mockGithubGateway := &rivalMockGithubGateway{
		GetUserFunc: func(ctx context.Context, username string) (*gateway.GithubUser, error) {
			return nil, errors.New("not found")
		},
	}

	usecase := NewRivalUsecase(mockRivalRepo, mockGithubGateway)

	rival, err := usecase.AddRival(ctx, 1, "nonexistentuser")

	assert.Error(t, err)
	assert.Nil(t, rival)
	assert.Contains(t, err.Error(), "見つかりません")
}

func TestAddRival_AlreadyExists(t *testing.T) {
	ctx := context.Background()
	mockRivalRepo := &rivalMockRivalRepository{
		CountByUserIDFunc: func(ctx context.Context, userID uint64) (int64, error) {
			return 0, nil
		},
		ExistsByUserIDAndRivalGithubUserIDFunc: func(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error) {
			return true, nil
		},
	}
	mockGithubGateway := &rivalMockGithubGateway{
		GetUserFunc: func(ctx context.Context, username string) (*gateway.GithubUser, error) {
			return &gateway.GithubUser{
				ID:    12345,
				Login: "rivaluser",
			}, nil
		},
	}

	usecase := NewRivalUsecase(mockRivalRepo, mockGithubGateway)

	rival, err := usecase.AddRival(ctx, 1, "rivaluser")

	assert.Error(t, err)
	assert.Nil(t, rival)
	assert.Contains(t, err.Error(), "既にライバルとして登録されています")
}

func TestRemoveRival_Success(t *testing.T) {
	ctx := context.Background()
	mockRivalRepo := &rivalMockRivalRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Rival, error) {
			return &models.Rival{
				ID:     1,
				UserID: 1,
			}, nil
		},
		DeleteFunc: func(ctx context.Context, id uint64) error {
			return nil
		},
	}
	mockGithubGateway := &rivalMockGithubGateway{}

	usecase := NewRivalUsecase(mockRivalRepo, mockGithubGateway)

	err := usecase.RemoveRival(ctx, 1, 1)

	assert.NoError(t, err)
}

func TestRemoveRival_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRivalRepo := &rivalMockRivalRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Rival, error) {
			return nil, errors.New("not found")
		},
	}
	mockGithubGateway := &rivalMockGithubGateway{}

	usecase := NewRivalUsecase(mockRivalRepo, mockGithubGateway)

	err := usecase.RemoveRival(ctx, 1, 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "見つかりません")
}

func TestRemoveRival_Unauthorized(t *testing.T) {
	ctx := context.Background()
	mockRivalRepo := &rivalMockRivalRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Rival, error) {
			return &models.Rival{
				ID:     1,
				UserID: 2, // Different user
			}, nil
		},
	}
	mockGithubGateway := &rivalMockGithubGateway{}

	usecase := NewRivalUsecase(mockRivalRepo, mockGithubGateway)

	err := usecase.RemoveRival(ctx, 1, 1) // User 1 trying to delete User 2's rival

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "権限がありません")
}
