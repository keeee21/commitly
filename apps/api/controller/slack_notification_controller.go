package controller

import (
	"net/http"
	"strings"

	"github.com/keeee21/commitly/api/dto"
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
// @Summary      Slack通知設定を取得
// @Description  現在のSlack通知設定を返す（Webhook URLはマスク済み）
// @Tags         notifications
// @Produce      json
// @Success      200 {object} dto.SlackNotificationSettingResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/notifications/slack [get]
func (ctrl *slackNotificationController) GetSetting(c echo.Context) error {
	user := c.Get("user").(*models.User)

	setting, err := ctrl.slackNotificationUsecase.GetSetting(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Slack通知設定の取得に失敗しました",
		})
	}

	if setting == nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Slack通知設定が見つかりません",
		})
	}

	maskedURL := maskWebhookURL(setting.WebhookURL)

	return c.JSON(http.StatusOK, dto.SlackNotificationSettingResponse{
		ID:         setting.ID,
		WebhookURL: maskedURL,
		IsEnabled:  setting.IsEnabled,
		CreatedAt:  setting.CreatedAt,
		UpdatedAt:  setting.UpdatedAt,
	})
}

// maskWebhookURL Webhook URLをマスク
func maskWebhookURL(url string) string {
	if len(url) <= 40 {
		return url[:20] + "..."
	}
	return url[:40] + "..."
}

// Create Slack通知設定を作成
// @Summary      Slack通知設定を作成
// @Description  Slack Webhook URLを登録して通知設定を作成する
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateSlackNotificationRequest true "Slack通知設定作成リクエスト"
// @Success      201 {object} dto.SlackNotificationSettingResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/notifications/slack [post]
func (ctrl *slackNotificationController) Create(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var req dto.CreateSlackNotificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "リクエストが不正です",
		})
	}

	if req.WebhookURL == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Webhook URLを入力してください",
		})
	}

	if !strings.HasPrefix(req.WebhookURL, slackWebhookURLPrefix) {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "無効なSlack Webhook URLです。https://hooks.slack.com/services/ で始まるURLを入力してください",
		})
	}

	setting, err := ctrl.slackNotificationUsecase.Create(c.Request().Context(), user.ID, req.WebhookURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Slack通知設定の作成に失敗しました",
		})
	}

	maskedURL := maskWebhookURL(setting.WebhookURL)

	return c.JSON(http.StatusCreated, dto.SlackNotificationSettingResponse{
		ID:         setting.ID,
		WebhookURL: maskedURL,
		IsEnabled:  setting.IsEnabled,
		CreatedAt:  setting.CreatedAt,
		UpdatedAt:  setting.UpdatedAt,
	})
}

// UpdateEnabled Slack通知の有効/無効を更新
// @Summary      Slack通知の有効/無効を更新
// @Description  Slack通知設定の有効/無効を切り替える
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        request body dto.UpdateEnabledRequest true "有効/無効更新リクエスト"
// @Success      200 {object} dto.UpdateEnabledResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/notifications/slack [put]
func (ctrl *slackNotificationController) UpdateEnabled(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var req dto.UpdateEnabledRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "リクエストが不正です",
		})
	}

	if err := ctrl.slackNotificationUsecase.UpdateEnabled(c.Request().Context(), user.ID, req.IsEnabled); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Slack通知設定の更新に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, dto.UpdateEnabledResponse{
		IsEnabled: req.IsEnabled,
	})
}

// Delete Slack通知設定を削除
// @Summary      Slack通知設定を削除
// @Description  Slack通知設定を削除する
// @Tags         notifications
// @Success      204
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/notifications/slack [delete]
func (ctrl *slackNotificationController) Delete(c echo.Context) error {
	user := c.Get("user").(*models.User)

	if err := ctrl.slackNotificationUsecase.Delete(c.Request().Context(), user.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Slack通知設定の削除に失敗しました",
		})
	}

	return c.NoContent(http.StatusNoContent)
}
