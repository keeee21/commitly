package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/tests/mocks"
	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetWeeklyDashboard_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/weekly", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{
		ID:             1,
		GithubUserID:   12345,
		GithubUsername: "testuser",
	}
	c.Set("user", user)

	rivals := []models.Rival{
		{ID: 1, RivalGithubUserID: 200, RivalGithubUsername: "rival1"},
	}

	dashboardData := &usecase.DashboardData{
		Period:    "weekly",
		StartDate: "2025-12-08",
		EndDate:   "2025-12-14",
		MyStats: usecase.UserCommitStats{
			GithubUserID:   12345,
			GithubUsername: "testuser",
			TotalCommits:   10,
		},
		Rivals: []usecase.UserCommitStats{
			{GithubUserID: 200, GithubUsername: "rival1", TotalCommits: 15},
		},
	}

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return rivals, nil
		},
	}

	mockDashboardUsecase := &mocks.MockDashboardUsecase{
		GetWeeklyDashboardFunc: func(ctx context.Context, user *models.User, rivals []models.Rival) (*usecase.DashboardData, error) {
			return dashboardData, nil
		},
	}

	ctrl := NewDashboardController(mockDashboardUsecase, mockRivalUsecase)
	err := ctrl.GetWeeklyDashboard(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"period":"weekly"`)
	assert.Contains(t, rec.Body.String(), `"total_commits":10`)
}

func TestGetWeeklyDashboard_RivalError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/weekly", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return nil, errors.New("database error")
		},
	}

	mockDashboardUsecase := &mocks.MockDashboardUsecase{}

	ctrl := NewDashboardController(mockDashboardUsecase, mockRivalUsecase)
	err := ctrl.GetWeeklyDashboard(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "ライバル情報の取得に失敗しました")
}

func TestGetWeeklyDashboard_DashboardError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/weekly", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return []models.Rival{}, nil
		},
	}

	mockDashboardUsecase := &mocks.MockDashboardUsecase{
		GetWeeklyDashboardFunc: func(ctx context.Context, user *models.User, rivals []models.Rival) (*usecase.DashboardData, error) {
			return nil, errors.New("dashboard error")
		},
	}

	ctrl := NewDashboardController(mockDashboardUsecase, mockRivalUsecase)
	err := ctrl.GetWeeklyDashboard(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "ダッシュボードデータの取得に失敗しました")
}

func TestGetMonthlyDashboard_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/monthly", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{
		ID:             1,
		GithubUserID:   12345,
		GithubUsername: "testuser",
	}
	c.Set("user", user)

	dashboardData := &usecase.DashboardData{
		Period:    "monthly",
		StartDate: "2025-12-01",
		EndDate:   "2025-12-14",
		MyStats: usecase.UserCommitStats{
			GithubUserID:   12345,
			GithubUsername: "testuser",
			TotalCommits:   50,
		},
		Rivals: []usecase.UserCommitStats{},
	}

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return []models.Rival{}, nil
		},
	}

	mockDashboardUsecase := &mocks.MockDashboardUsecase{
		GetMonthlyDashboardFunc: func(ctx context.Context, user *models.User, rivals []models.Rival) (*usecase.DashboardData, error) {
			return dashboardData, nil
		},
	}

	ctrl := NewDashboardController(mockDashboardUsecase, mockRivalUsecase)
	err := ctrl.GetMonthlyDashboard(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"period":"monthly"`)
	assert.Contains(t, rec.Body.String(), `"total_commits":50`)
}

func TestGetMonthlyDashboard_RivalError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/monthly", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return nil, errors.New("database error")
		},
	}

	mockDashboardUsecase := &mocks.MockDashboardUsecase{}

	ctrl := NewDashboardController(mockDashboardUsecase, mockRivalUsecase)
	err := ctrl.GetMonthlyDashboard(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
