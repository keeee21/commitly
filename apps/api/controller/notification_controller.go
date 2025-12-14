package controller

import (
	"net/http"
	"strconv"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
)

// INotificationController 通知コントローラーのインターフェース
type INotificationController interface {
	GetSettings(c echo.Context) error
	CreateSetting(c echo.Context) error
	UpdateSetting(c echo.Context) error
	DeleteSetting(c echo.Context) error
}

type notificationController struct {
	notificationUsecase usecase.INotificationUsecase
}

// NewNotificationController コンストラクタ
func NewNotificationController(notificationUsecase usecase.INotificationUsecase) INotificationController {
	return &notificationController{
		notificationUsecase: notificationUsecase,
	}
}

// CreateSettingRequest 通知設定作成リクエスト
type CreateSettingRequest struct {
	ChannelType string `json:"channel_type" validate:"required"`
	WebhookURL  string `json:"webhook_url"`
	LINEUserID  string `json:"line_user_id"`
}

// UpdateSettingRequest 通知設定更新リクエスト
type UpdateSettingRequest struct {
	IsEnabled  bool   `json:"is_enabled"`
	WebhookURL string `json:"webhook_url"`
	LINEUserID string `json:"line_user_id"`
}

// GetSettings 通知設定一覧を取得
func (ctrl *notificationController) GetSettings(c echo.Context) error {
	user := c.Get("user").(*models.User)

	settings, err := ctrl.notificationUsecase.GetSettings(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "通知設定の取得に失敗しました",
		})
	}

	// レスポンス用に変換
	response := make([]map[string]interface{}, len(settings))
	for i, setting := range settings {
		response[i] = map[string]interface{}{
			"id":           setting.ID,
			"channel_type": setting.ChannelType,
			"webhook_url":  setting.WebhookURL,
			"line_user_id": setting.LINEUserID,
			"is_enabled":   setting.IsEnabled,
			"created_at":   setting.CreatedAt,
			"updated_at":   setting.UpdatedAt,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"settings": response,
	})
}

// CreateSetting 通知設定を作成
func (ctrl *notificationController) CreateSetting(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var req CreateSettingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストが不正です",
		})
	}

	// チャンネルタイプのバリデーション
	channelType := models.ChannelType(req.ChannelType)
	if channelType != models.ChannelTypeLINE &&
		channelType != models.ChannelTypeSlack &&
		channelType != models.ChannelTypeDiscord {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "チャンネルタイプが不正です（line, slack, discord のいずれかを指定してください）",
		})
	}

	setting, err := ctrl.notificationUsecase.CreateSetting(
		c.Request().Context(),
		user.ID,
		channelType,
		req.WebhookURL,
		req.LINEUserID,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "通知設定の作成に失敗しました",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":           setting.ID,
		"channel_type": setting.ChannelType,
		"webhook_url":  setting.WebhookURL,
		"line_user_id": setting.LINEUserID,
		"is_enabled":   setting.IsEnabled,
		"created_at":   setting.CreatedAt,
	})
}

// UpdateSetting 通知設定を更新
func (ctrl *notificationController) UpdateSetting(c echo.Context) error {
	user := c.Get("user").(*models.User)

	settingIDStr := c.Param("id")
	settingID, err := strconv.ParseUint(settingIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "設定IDが不正です",
		})
	}

	var req UpdateSettingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストが不正です",
		})
	}

	setting, err := ctrl.notificationUsecase.UpdateSetting(
		c.Request().Context(),
		user.ID,
		settingID,
		req.IsEnabled,
		req.WebhookURL,
		req.LINEUserID,
	)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":           setting.ID,
		"channel_type": setting.ChannelType,
		"webhook_url":  setting.WebhookURL,
		"line_user_id": setting.LINEUserID,
		"is_enabled":   setting.IsEnabled,
		"updated_at":   setting.UpdatedAt,
	})
}

// DeleteSetting 通知設定を削除
func (ctrl *notificationController) DeleteSetting(c echo.Context) error {
	user := c.Get("user").(*models.User)

	settingIDStr := c.Param("id")
	settingID, err := strconv.ParseUint(settingIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "設定IDが不正です",
		})
	}

	if err := ctrl.notificationUsecase.DeleteSetting(c.Request().Context(), user.ID, settingID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
