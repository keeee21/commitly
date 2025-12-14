package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/tests/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupRivalControllerTest() (*echo.Echo, *models.User) {
	e := echo.New()
	user := &models.User{
		ID:             1,
		GithubUserID:   12345,
		GithubUsername: "testuser",
	}
	return e, user
}

func TestGetRivals_Success(t *testing.T) {
	e, user := setupRivalControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/rivals", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return []models.Rival{
				{
					ID:                  1,
					UserID:              userID,
					RivalGithubUserID:   100,
					RivalGithubUsername: "rival1",
					RivalAvatarURL:      "https://avatar1.url",
					CreatedAt:           time.Now(),
				},
			}, nil
		},
	}

	ctrl := NewRivalController(mockRivalUsecase)
	err := ctrl.GetRivals(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"count":1`)
	assert.Contains(t, rec.Body.String(), `"github_username":"rival1"`)
}

func TestGetRivals_Empty(t *testing.T) {
	e, user := setupRivalControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/rivals", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return []models.Rival{}, nil
		},
	}

	ctrl := NewRivalController(mockRivalUsecase)
	err := ctrl.GetRivals(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"count":0`)
}

func TestGetRivals_Error(t *testing.T) {
	e, user := setupRivalControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/rivals", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		GetRivalsFunc: func(ctx context.Context, userID uint64) ([]models.Rival, error) {
			return nil, errors.New("database error")
		},
	}

	ctrl := NewRivalController(mockRivalUsecase)
	err := ctrl.GetRivals(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAddRival_Success(t *testing.T) {
	e, user := setupRivalControllerTest()
	body := `{"username": "newrival"}`
	req := httptest.NewRequest(http.MethodPost, "/api/rivals", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		AddRivalFunc: func(ctx context.Context, userID uint64, rivalUsername string) (*models.Rival, error) {
			return &models.Rival{
				ID:                  1,
				UserID:              userID,
				RivalGithubUserID:   200,
				RivalGithubUsername: rivalUsername,
				RivalAvatarURL:      "https://avatar.url",
				CreatedAt:           time.Now(),
			}, nil
		},
	}

	ctrl := NewRivalController(mockRivalUsecase)
	err := ctrl.AddRival(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"github_username":"newrival"`)
}

func TestAddRival_EmptyUsername(t *testing.T) {
	e, user := setupRivalControllerTest()
	body := `{"username": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/rivals", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{}

	ctrl := NewRivalController(mockRivalUsecase)
	err := ctrl.AddRival(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "ユーザー名を指定してください")
}

func TestAddRival_InvalidRequest(t *testing.T) {
	e, user := setupRivalControllerTest()
	body := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/rivals", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{}

	ctrl := NewRivalController(mockRivalUsecase)
	err := ctrl.AddRival(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddRival_UsecaseError(t *testing.T) {
	e, user := setupRivalControllerTest()
	body := `{"username": "newrival"}`
	req := httptest.NewRequest(http.MethodPost, "/api/rivals", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		AddRivalFunc: func(ctx context.Context, userID uint64, rivalUsername string) (*models.Rival, error) {
			return nil, errors.New("ライバル登録数が上限に達しています")
		},
	}

	ctrl := NewRivalController(mockRivalUsecase)
	err := ctrl.AddRival(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "上限")
}

func TestRemoveRival_Success(t *testing.T) {
	e, user := setupRivalControllerTest()
	req := httptest.NewRequest(http.MethodDelete, "/api/rivals/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		RemoveRivalFunc: func(ctx context.Context, userID uint64, rivalID uint64) error {
			return nil
		},
	}

	ctrl := NewRivalController(mockRivalUsecase)
	err := ctrl.RemoveRival(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestRemoveRival_InvalidID(t *testing.T) {
	e, user := setupRivalControllerTest()
	req := httptest.NewRequest(http.MethodDelete, "/api/rivals/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{}

	ctrl := NewRivalController(mockRivalUsecase)
	err := ctrl.RemoveRival(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "ライバルIDが不正です")
}

func TestRemoveRival_NotFound(t *testing.T) {
	e, user := setupRivalControllerTest()
	req := httptest.NewRequest(http.MethodDelete, "/api/rivals/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")
	c.Set("user", user)

	mockRivalUsecase := &mocks.MockRivalUsecase{
		RemoveRivalFunc: func(ctx context.Context, userID uint64, rivalID uint64) error {
			return errors.New("ライバルが見つかりません")
		},
	}

	ctrl := NewRivalController(mockRivalUsecase)
	err := ctrl.RemoveRival(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
