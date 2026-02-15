package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/keeee21/commitly/api/dto"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/tests/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupActivityControllerTest() (*echo.Echo, *models.User) {
	e := echo.New()
	user := &models.User{
		ID:             1,
		GithubUserID:   12345,
		GithubUsername: "testuser",
		AvatarURL:      "https://avatar.url",
	}
	return e, user
}

func TestGetActivityStream_Success(t *testing.T) {
	e, user := setupActivityControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/activity/stream", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return []models.Rival{}, nil
		},
	}
	mockActivityUsecase := &mocks.MockActivityUsecase{
		GetActivityStreamFunc: func(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.ActivityStreamResponse, error) {
			return &dto.ActivityStreamResponse{
				Activities: []dto.ActivityItem{
					{GithubUsername: "testuser", Repository: "my-repo", CommitCount: 3, Date: "2026-02-15"},
				},
			}, nil
		},
	}

	ctrl := NewActivityController(mockActivityUsecase, mockRivalUsecase)
	err := ctrl.GetActivityStream(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"github_username":"testuser"`)
	assert.Contains(t, rec.Body.String(), `"repository":"my-repo"`)
}

func TestGetActivityStream_RivalError(t *testing.T) {
	e, user := setupActivityControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/activity/stream", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return nil, errors.New("database error")
		},
	}
	mockActivityUsecase := &mocks.MockActivityUsecase{}

	ctrl := NewActivityController(mockActivityUsecase, mockRivalUsecase)
	err := ctrl.GetActivityStream(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "ライバル情報の取得に失敗しました")
}

func TestGetActivityStream_UsecaseError(t *testing.T) {
	e, user := setupActivityControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/activity/stream", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return []models.Rival{}, nil
		},
	}
	mockActivityUsecase := &mocks.MockActivityUsecase{
		GetActivityStreamFunc: func(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.ActivityStreamResponse, error) {
			return nil, errors.New("データ取得エラー")
		},
	}

	ctrl := NewActivityController(mockActivityUsecase, mockRivalUsecase)
	err := ctrl.GetActivityStream(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "アクティビティデータの取得に失敗しました")
}

func TestGetRhythm_Success(t *testing.T) {
	e, user := setupActivityControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/activity/rhythm", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return []models.Rival{}, nil
		},
	}
	mockActivityUsecase := &mocks.MockActivityUsecase{
		GetRhythmFunc: func(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.RhythmResponse, error) {
			return &dto.RhythmResponse{
				Users: []dto.UserRhythm{
					{
						GithubUsername: "testuser",
						PatternLabel:   "安定型",
						WeeklyRhythm: dto.WeeklyRhythm{
							Mon: true, Tue: true, Wed: true, Thu: true, Fri: true,
						},
					},
				},
				Period: "2026-02-09/2026-02-15",
			}, nil
		},
	}

	ctrl := NewActivityController(mockActivityUsecase, mockRivalUsecase)
	err := ctrl.GetRhythm(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"pattern_label":"安定型"`)
	assert.Contains(t, rec.Body.String(), `"period":"2026-02-09/2026-02-15"`)
}

func TestGetRhythm_RivalError(t *testing.T) {
	e, user := setupActivityControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/activity/rhythm", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return nil, errors.New("database error")
		},
	}
	mockActivityUsecase := &mocks.MockActivityUsecase{}

	ctrl := NewActivityController(mockActivityUsecase, mockRivalUsecase)
	err := ctrl.GetRhythm(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetRhythm_UsecaseError(t *testing.T) {
	e, user := setupActivityControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/activity/rhythm", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return []models.Rival{}, nil
		},
	}
	mockActivityUsecase := &mocks.MockActivityUsecase{
		GetRhythmFunc: func(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.RhythmResponse, error) {
			return nil, errors.New("データ取得エラー")
		},
	}

	ctrl := NewActivityController(mockActivityUsecase, mockRivalUsecase)
	err := ctrl.GetRhythm(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "リズムデータの取得に失敗しました")
}
