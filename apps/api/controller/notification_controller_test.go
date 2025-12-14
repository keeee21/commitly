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

func TestGetSettings_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/notifications/settings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	settings := []models.NotificationSetting{
		{ID: 1, UserID: 1, ChannelType: models.ChannelTypeLINE, IsEnabled: true},
		{ID: 2, UserID: 1, ChannelType: models.ChannelTypeSlack, WebhookURL: "https://hooks.slack.com/xxx", IsEnabled: true},
	}

	mockUsecase := &mocks.MockNotificationUsecase{
		GetSettingsFunc: func(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
			return settings, nil
		},
	}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.GetSettings(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"settings"`)
	assert.Contains(t, rec.Body.String(), `"channel_type":"line"`)
	assert.Contains(t, rec.Body.String(), `"channel_type":"slack"`)
}

func TestGetSettings_Empty(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/notifications/settings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{
		GetSettingsFunc: func(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
			return []models.NotificationSetting{}, nil
		},
	}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.GetSettings(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"settings":[]`)
}

func TestGetSettings_Error(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/notifications/settings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{
		GetSettingsFunc: func(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
			return nil, errors.New("database error")
		},
	}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.GetSettings(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "通知設定の取得に失敗しました")
}

func TestCreateSetting_Success(t *testing.T) {
	e := echo.New()
	reqBody := `{"channel_type":"slack","webhook_url":"https://hooks.slack.com/xxx"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notifications/settings", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{
		CreateSettingFunc: func(ctx context.Context, userID uint64, channelType models.ChannelType, webhookURL, lineUserID string) (*models.NotificationSetting, error) {
			return &models.NotificationSetting{
				ID:          1,
				UserID:      userID,
				ChannelType: channelType,
				WebhookURL:  webhookURL,
				IsEnabled:   true,
			}, nil
		},
	}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.CreateSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"channel_type":"slack"`)
	assert.Contains(t, rec.Body.String(), `"is_enabled":true`)
}

func TestCreateSetting_LINE(t *testing.T) {
	e := echo.New()
	reqBody := `{"channel_type":"line","line_user_id":"U1234567890"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notifications/settings", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{
		CreateSettingFunc: func(ctx context.Context, userID uint64, channelType models.ChannelType, webhookURL, lineUserID string) (*models.NotificationSetting, error) {
			return &models.NotificationSetting{
				ID:          1,
				UserID:      userID,
				ChannelType: channelType,
				LINEUserID:  lineUserID,
				IsEnabled:   true,
			}, nil
		},
	}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.CreateSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"channel_type":"line"`)
}

func TestCreateSetting_InvalidChannelType(t *testing.T) {
	e := echo.New()
	reqBody := `{"channel_type":"invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notifications/settings", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.CreateSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "チャンネルタイプが不正です")
}

func TestCreateSetting_InvalidJSON(t *testing.T) {
	e := echo.New()
	reqBody := `{"channel_type":}`
	req := httptest.NewRequest(http.MethodPost, "/api/notifications/settings", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.CreateSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "リクエストが不正です")
}

func TestCreateSetting_UsecaseError(t *testing.T) {
	e := echo.New()
	reqBody := `{"channel_type":"slack","webhook_url":"https://hooks.slack.com/xxx"}`
	req := httptest.NewRequest(http.MethodPost, "/api/notifications/settings", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{
		CreateSettingFunc: func(ctx context.Context, userID uint64, channelType models.ChannelType, webhookURL, lineUserID string) (*models.NotificationSetting, error) {
			return nil, errors.New("database error")
		},
	}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.CreateSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "通知設定の作成に失敗しました")
}

func TestUpdateSetting_Success(t *testing.T) {
	e := echo.New()
	reqBody := `{"is_enabled":false,"webhook_url":"https://new.url"}`
	req := httptest.NewRequest(http.MethodPut, "/api/notifications/settings/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{
		UpdateSettingFunc: func(ctx context.Context, userID uint64, settingID uint64, isEnabled bool, webhookURL, lineUserID string) (*models.NotificationSetting, error) {
			return &models.NotificationSetting{
				ID:          settingID,
				UserID:      userID,
				ChannelType: models.ChannelTypeSlack,
				WebhookURL:  webhookURL,
				IsEnabled:   isEnabled,
			}, nil
		},
	}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.UpdateSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"is_enabled":false`)
	assert.Contains(t, rec.Body.String(), `"webhook_url":"https://new.url"`)
}

func TestUpdateSetting_InvalidID(t *testing.T) {
	e := echo.New()
	reqBody := `{"is_enabled":false}`
	req := httptest.NewRequest(http.MethodPut, "/api/notifications/settings/invalid", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.UpdateSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "設定IDが不正です")
}

func TestUpdateSetting_InvalidJSON(t *testing.T) {
	e := echo.New()
	reqBody := `{"is_enabled":}`
	req := httptest.NewRequest(http.MethodPut, "/api/notifications/settings/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.UpdateSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "リクエストが不正です")
}

func TestUpdateSetting_UsecaseError(t *testing.T) {
	e := echo.New()
	reqBody := `{"is_enabled":false}`
	req := httptest.NewRequest(http.MethodPut, "/api/notifications/settings/1", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{
		UpdateSettingFunc: func(ctx context.Context, userID uint64, settingID uint64, isEnabled bool, webhookURL, lineUserID string) (*models.NotificationSetting, error) {
			return nil, errors.New("設定が見つかりません")
		},
	}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.UpdateSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "設定が見つかりません")
}

func TestDeleteSetting_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/notifications/settings/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{
		DeleteSettingFunc: func(ctx context.Context, userID uint64, settingID uint64) error {
			return nil
		},
	}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.DeleteSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteSetting_InvalidID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/notifications/settings/invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("invalid")

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.DeleteSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "設定IDが不正です")
}

func TestDeleteSetting_UsecaseError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/notifications/settings/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	user := &models.User{ID: 1}
	c.Set("user", user)

	mockUsecase := &mocks.MockNotificationUsecase{
		DeleteSettingFunc: func(ctx context.Context, userID uint64, settingID uint64) error {
			return errors.New("設定が見つかりません")
		},
	}

	ctrl := NewNotificationController(mockUsecase)
	err := ctrl.DeleteSetting(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "設定が見つかりません")
}
