package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/tests/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-GitHub-User-ID", "12345")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedUser := &models.User{
		ID:             1,
		GithubUserID:   12345,
		GithubUsername: "testuser",
	}

	mockUserUsecase := &mocks.MockUserUsecase{
		GetUserByGithubUserIDFunc: func(ctx context.Context, githubUserID uint64) (*models.User, error) {
			return expectedUser, nil
		},
	}

	handler := AuthMiddleware(mockUserUsecase)(func(c echo.Context) error {
		user := c.Get("user").(*models.User)
		assert.Equal(t, expectedUser, user)
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// No X-GitHub-User-ID header
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserUsecase := &mocks.MockUserUsecase{}

	handler := AuthMiddleware(mockUserUsecase)(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "認証が必要です")
}

func TestAuthMiddleware_InvalidHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-GitHub-User-ID", "invalid")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserUsecase := &mocks.MockUserUsecase{}

	handler := AuthMiddleware(mockUserUsecase)(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "認証情報が不正です")
}

func TestAuthMiddleware_UserNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-GitHub-User-ID", "99999")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUserUsecase := &mocks.MockUserUsecase{
		GetUserByGithubUserIDFunc: func(ctx context.Context, githubUserID uint64) (*models.User, error) {
			return nil, errors.New("not found")
		},
	}

	handler := AuthMiddleware(mockUserUsecase)(func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	err := handler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "ユーザーが見つかりません")
}
