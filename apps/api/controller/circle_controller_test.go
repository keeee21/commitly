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

func setupCircleControllerTest() (*echo.Echo, *models.User) {
	e := echo.New()
	user := &models.User{
		ID:             1,
		GithubUserID:   12345,
		GithubUsername: "testuser",
		AvatarURL:      "https://avatar.url",
	}
	return e, user
}

func TestCircleGetCircles_Success(t *testing.T) {
	e, user := setupCircleControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/circles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{
		GetCirclesFunc: func(ctx context.Context, userID uint64) ([]models.Circle, error) {
			return []models.Circle{
				{
					ID: 1, Name: "Circle1", OwnerUserID: userID, InviteCode: "abcd1234",
					CreatedAt: time.Now(),
					Members: []models.CircleMember{
						{ID: 1, CircleID: 1, UserID: userID, User: models.User{GithubUsername: "testuser", AvatarURL: "https://avatar.url"}},
					},
				},
			}, nil
		},
	}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.GetCircles(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"count":1`)
	assert.Contains(t, rec.Body.String(), `"name":"Circle1"`)
}

func TestCircleGetCircles_Error(t *testing.T) {
	e, user := setupCircleControllerTest()
	req := httptest.NewRequest(http.MethodGet, "/api/circles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{
		GetCirclesFunc: func(ctx context.Context, userID uint64) ([]models.Circle, error) {
			return nil, errors.New("database error")
		},
	}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.GetCircles(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestCircleCreateCircle_Success(t *testing.T) {
	e, user := setupCircleControllerTest()
	body := `{"name": "新サークル"}`
	req := httptest.NewRequest(http.MethodPost, "/api/circles", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{
		CreateCircleFunc: func(ctx context.Context, userID uint64, name string) (*models.Circle, error) {
			return &models.Circle{
				ID: 1, Name: name, OwnerUserID: userID, InviteCode: "abcd1234",
				CreatedAt: time.Now(),
				Members: []models.CircleMember{
					{ID: 1, CircleID: 1, UserID: userID, User: models.User{GithubUsername: "testuser", AvatarURL: "https://avatar.url"}},
				},
			}, nil
		},
	}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.CreateCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"name":"新サークル"`)
	assert.Contains(t, rec.Body.String(), `"invite_code":"abcd1234"`)
}

func TestCircleCreateCircle_EmptyName(t *testing.T) {
	e, user := setupCircleControllerTest()
	body := `{"name": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/circles", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.CreateCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "サークル名を入力してください")
}

func TestCircleCreateCircle_Error(t *testing.T) {
	e, user := setupCircleControllerTest()
	body := `{"name": "テスト"}`
	req := httptest.NewRequest(http.MethodPost, "/api/circles", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{
		CreateCircleFunc: func(ctx context.Context, userID uint64, name string) (*models.Circle, error) {
			return nil, errors.New("サークル作成数が上限に達しています")
		},
	}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.CreateCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "上限")
}

func TestCircleJoinCircle_Success(t *testing.T) {
	e, user := setupCircleControllerTest()
	body := `{"invite_code": "abcd1234"}`
	req := httptest.NewRequest(http.MethodPost, "/api/circles/join", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{
		JoinCircleFunc: func(ctx context.Context, userID uint64, inviteCode string) (*models.Circle, error) {
			return &models.Circle{
				ID: 1, Name: "Circle1", OwnerUserID: 2, InviteCode: inviteCode,
				CreatedAt: time.Now(),
				Members: []models.CircleMember{
					{ID: 1, CircleID: 1, UserID: 2, User: models.User{GithubUsername: "owner", AvatarURL: "https://avatar.url"}},
					{ID: 2, CircleID: 1, UserID: userID, User: models.User{GithubUsername: "testuser", AvatarURL: "https://avatar.url"}},
				},
			}, nil
		},
	}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.JoinCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"name":"Circle1"`)
}

func TestCircleJoinCircle_EmptyCode(t *testing.T) {
	e, user := setupCircleControllerTest()
	body := `{"invite_code": ""}`
	req := httptest.NewRequest(http.MethodPost, "/api/circles/join", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.JoinCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "招待コードを入力してください")
}

func TestCircleJoinCircle_Error(t *testing.T) {
	e, user := setupCircleControllerTest()
	body := `{"invite_code": "invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/api/circles/join", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{
		JoinCircleFunc: func(ctx context.Context, userID uint64, inviteCode string) (*models.Circle, error) {
			return nil, errors.New("招待コードが無効です")
		},
	}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.JoinCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "招待コードが無効です")
}

func TestCircleLeaveCircle_Success(t *testing.T) {
	e, user := setupCircleControllerTest()
	req := httptest.NewRequest(http.MethodDelete, "/api/circles/1/leave", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{
		LeaveCircleFunc: func(ctx context.Context, userID uint64, circleID uint64) error {
			return nil
		},
	}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.LeaveCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestCircleLeaveCircle_InvalidID(t *testing.T) {
	e, user := setupCircleControllerTest()
	req := httptest.NewRequest(http.MethodDelete, "/api/circles/invalid/leave", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.LeaveCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "サークルIDが不正です")
}

func TestCircleLeaveCircle_Error(t *testing.T) {
	e, user := setupCircleControllerTest()
	req := httptest.NewRequest(http.MethodDelete, "/api/circles/1/leave", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{
		LeaveCircleFunc: func(ctx context.Context, userID uint64, circleID uint64) error {
			return errors.New("オーナーはサークルを退会できません")
		},
	}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.LeaveCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "オーナー")
}

func TestCircleDeleteCircle_Success(t *testing.T) {
	e, user := setupCircleControllerTest()
	req := httptest.NewRequest(http.MethodDelete, "/api/circles/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{
		DeleteCircleFunc: func(ctx context.Context, userID uint64, circleID uint64) error {
			return nil
		},
	}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.DeleteCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestCircleDeleteCircle_InvalidID(t *testing.T) {
	e, user := setupCircleControllerTest()
	req := httptest.NewRequest(http.MethodDelete, "/api/circles/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.DeleteCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "サークルIDが不正です")
}

func TestCircleDeleteCircle_Error(t *testing.T) {
	e, user := setupCircleControllerTest()
	req := httptest.NewRequest(http.MethodDelete, "/api/circles/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user", user)

	mockUsecase := &mocks.MockCircleUsecase{
		DeleteCircleFunc: func(ctx context.Context, userID uint64, circleID uint64) error {
			return errors.New("サークルを削除する権限がありません")
		},
	}

	ctrl := NewCircleController(mockUsecase)
	err := ctrl.DeleteCircle(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "権限がありません")
}
