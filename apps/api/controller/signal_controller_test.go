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

func setupSignalControllerTest() (*echo.Echo, *models.User) {
	e := echo.New()
	user := &models.User{
		ID:             1,
		GithubUserID:   12345,
		GithubUsername: "testuser",
		AvatarURL:      "https://avatar.url",
	}
	return e, user
}

func TestGetSignals_Success(t *testing.T) {
	e, user := setupSignalControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/circles/1/signals", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user", user)

	mockUsecase := &mocks.MockSignalUsecase{
		GetSignalsFunc: func(ctx context.Context, userID uint64, circleID uint64) ([]usecase.Signal, error) {
			return []usecase.Signal{
				{Type: "same_day", Date: "2026-02-14", Usernames: []string{"tanaka"}, AvatarURLs: []string{"https://avatar/tanaka"}, Detail: "同じ日にコミット"},
			}, nil
		},
	}

	ctrl := NewSignalController(mockUsecase)
	err := ctrl.GetSignals(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"type":"same_day"`)
	assert.Contains(t, rec.Body.String(), `"tanaka"`)
}

func TestGetSignals_InvalidID(t *testing.T) {
	e, user := setupSignalControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/circles/abc/signals", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")
	c.Set("user", user)

	mockUsecase := &mocks.MockSignalUsecase{}
	ctrl := NewSignalController(mockUsecase)
	err := ctrl.GetSignals(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "サークルIDが不正です")
}

func TestGetSignals_Error(t *testing.T) {
	e, user := setupSignalControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/circles/1/signals", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user", user)

	mockUsecase := &mocks.MockSignalUsecase{
		GetSignalsFunc: func(ctx context.Context, userID uint64, circleID uint64) ([]usecase.Signal, error) {
			return nil, errors.New("このサークルのメンバーではありません")
		},
	}

	ctrl := NewSignalController(mockUsecase)
	err := ctrl.GetSignals(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "メンバーではありません")
}

func TestGetRecentSignals_Success(t *testing.T) {
	e, user := setupSignalControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/signals/recent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockUsecase := &mocks.MockSignalUsecase{
		GetRecentSignalsFunc: func(ctx context.Context, userID uint64) ([]usecase.Signal, error) {
			return []usecase.Signal{
				{Type: "same_day", Date: "2026-02-14", Usernames: []string{"tanaka"}, AvatarURLs: []string{"https://avatar/tanaka"}, Detail: "同じ日にコミット", CircleID: 1, CircleName: "テスト"},
			}, nil
		},
	}

	ctrl := NewSignalController(mockUsecase)
	err := ctrl.GetRecentSignals(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"circle_id":1`)
	assert.Contains(t, rec.Body.String(), `"circle_name":"テスト"`)
}

func TestGetRecentSignals_Error(t *testing.T) {
	e, user := setupSignalControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/signals/recent", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockUsecase := &mocks.MockSignalUsecase{
		GetRecentSignalsFunc: func(ctx context.Context, userID uint64) ([]usecase.Signal, error) {
			return nil, errors.New("something went wrong")
		},
	}

	ctrl := NewSignalController(mockUsecase)
	err := ctrl.GetRecentSignals(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "シグナルの取得に失敗しました")
}
