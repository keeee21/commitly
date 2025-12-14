package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/tests/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCallback_Success(t *testing.T) {
	e := echo.New()
	body := `{"github_user_id": 12345, "github_username": "testuser", "email": "test@example.com", "avatar_url": "https://avatar.url"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/callback", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserUsecase := &mocks.MockUserUsecase{
		GetOrCreateUserFunc: func(ctx context.Context, githubUserID uint64, githubUsername, email, avatarURL string) (*models.User, error) {
			return &models.User{
				ID:             1,
				GithubUserID:   githubUserID,
				GithubUsername: githubUsername,
				Email:          email,
				AvatarURL:      avatarURL,
			}, nil
		},
	}

	ctrl := NewAuthController(mockUserUsecase)
	err := ctrl.Callback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"github_user_id":12345`)
	assert.Contains(t, rec.Body.String(), `"github_username":"testuser"`)
	// email should NOT be in response
	assert.NotContains(t, rec.Body.String(), `"email"`)
}

func TestCallback_InvalidRequest(t *testing.T) {
	e := echo.New()
	body := `invalid json`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/callback", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserUsecase := &mocks.MockUserUsecase{}

	ctrl := NewAuthController(mockUserUsecase)
	err := ctrl.Callback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "リクエストが不正です")
}

func TestCallback_MissingRequiredFields(t *testing.T) {
	e := echo.New()
	body := `{"github_user_id": 0, "github_username": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/callback", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserUsecase := &mocks.MockUserUsecase{}

	ctrl := NewAuthController(mockUserUsecase)
	err := ctrl.Callback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Githubユーザー情報が不正です")
}

func TestCallback_UsecaseError(t *testing.T) {
	e := echo.New()
	body := `{"github_user_id": 12345, "github_username": "testuser"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/callback", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserUsecase := &mocks.MockUserUsecase{
		GetOrCreateUserFunc: func(ctx context.Context, githubUserID uint64, githubUsername, email, avatarURL string) (*models.User, error) {
			return nil, errors.New("database error")
		},
	}

	ctrl := NewAuthController(mockUserUsecase)
	err := ctrl.Callback(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "ユーザー情報の保存に失敗しました")
}

func TestLogout_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserUsecase := &mocks.MockUserUsecase{}

	ctrl := NewAuthController(mockUserUsecase)
	err := ctrl.Logout(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "ログアウトしました")
}
