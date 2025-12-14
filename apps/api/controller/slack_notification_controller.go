package controller

import (
	"net/http"
	"strings"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
)

const slackWebhookURLPrefix = "https://hooks.slack.com/services/"

// ISlackNotificationController Slack通知コントローラーのインターフェース
type ISlackNotificationController interface {
	GetSetting(c echo.Context) error
	Create(c echo.Context) error
	UpdateEnabled(c echo.Context) error
	Delete(c echo.Context) error
}

type slackNotificationController struct {
	slackNotificationUsecase usecase.ISlackNotificationUsecase
}

// NewSlackNotificationController コンストラクタ
func NewSlackNotificationController(slackNotificationUsecase usecase.ISlackNotificationUsecase) ISlackNotificationController {
	return &slackNotificationController{
		slackNotificationUsecase: slackNotificationUsecase,
	}
}

// GetSetting Slack通知設定を取得
func (ctrl *slackNotificationController) GetSetting(c echo.Context) error {
	user := c.Get("user").(*models.User)

	setting, err := ctrl.slackNotificationUsecase.GetSetting(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Slack通知設定の取得に失敗しました",
		})
	}

	if setting == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Slack通知設定が見つかりません",
		})
	}

	// Webhook URLの一部をマスク（セキュリティのため）
	maskedURL := maskWebhookURL(setting.WebhookURL)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":          setting.ID,
		"webhook_url": maskedURL,
		"is_enabled":  setting.IsEnabled,
		"created_at":  setting.CreatedAt,
		"updated_at":  setting.UpdatedAt,
	})
}

// maskWebhookURL Webhook URLをマスク
func maskWebhookURL(url string) string {
	if len(url) <= 40 {
		return url[:20] + "..."
	}
	return url[:40] + "..."
}

// CreateRequest Slack通知設定作成リクエスト
type CreateRequest struct {
	WebhookURL string `json:"webhook_url"`
}

// Create Slack通知設定を作成
func (ctrl *slackNotificationController) Create(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストが不正です",
		})
	}

	// Webhook URL のバリデーション
	if req.WebhookURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Webhook URLを入力してください",
		})
	}

	if !strings.HasPrefix(req.WebhookURL, slackWebhookURLPrefix) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なSlack Webhook URLです。https://hooks.slack.com/services/ で始まるURLを入力してください",
		})
	}

	setting, err := ctrl.slackNotificationUsecase.Create(c.Request().Context(), user.ID, req.WebhookURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Slack通知設定の作成に失敗しました",
		})
	}

	maskedURL := maskWebhookURL(setting.WebhookURL)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":          setting.ID,
		"webhook_url": maskedURL,
		"is_enabled":  setting.IsEnabled,
		"created_at":  setting.CreatedAt,
		"updated_at":  setting.UpdatedAt,
	})
}

// UpdateEnabledRequest 有効/無効更新リクエスト
type UpdateEnabledRequest struct {
	IsEnabled bool `json:"is_enabled"`
}

// UpdateEnabled Slack通知の有効/無効を更新
func (ctrl *slackNotificationController) UpdateEnabled(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var req UpdateEnabledRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストが不正です",
		})
	}

	if err := ctrl.slackNotificationUsecase.UpdateEnabled(c.Request().Context(), user.ID, req.IsEnabled); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Slack通知設定の更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"is_enabled": req.IsEnabled,
	})
}

// Delete Slack通知設定を削除
func (ctrl *slackNotificationController) Delete(c echo.Context) error {
	user := c.Get("user").(*models.User)

	if err := ctrl.slackNotificationUsecase.Delete(c.Request().Context(), user.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Slack通知設定の削除に失敗しました",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
