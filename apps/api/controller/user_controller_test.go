package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/tests/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetMe_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set user in context (simulating auth middleware)
	user := &models.User{
		ID:             1,
		GithubUserID:   12345,
		GithubUsername: "testuser",
		Email:          "test@example.com",
		AvatarURL:      "https://avatar.url",
		CreatedAt:      time.Now(),
	}
	c.Set("user", user)

	mockUserUsecase := &mocks.MockUserUsecase{}
	ctrl := NewUserController(mockUserUsecase)

	err := ctrl.GetMe(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"github_user_id":12345`)
	assert.Contains(t, rec.Body.String(), `"github_username":"testuser"`)
	assert.Contains(t, rec.Body.String(), `"avatar_url":"https://avatar.url"`)
	// email should NOT be in response
	assert.NotContains(t, rec.Body.String(), `"email"`)
}
