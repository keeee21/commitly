package controller

import (
	"net/http"
	"strconv"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
)

// IRivalController ライバルコントローラーのインターフェース
type IRivalController interface {
	GetRivals(c echo.Context) error
	AddRival(c echo.Context) error
	RemoveRival(c echo.Context) error
}

type rivalController struct {
	rivalUsecase usecase.IRivalUsecase
}

// NewRivalController コンストラクタ
func NewRivalController(rivalUsecase usecase.IRivalUsecase) IRivalController {
	return &rivalController{
		rivalUsecase: rivalUsecase,
	}
}

// AddRivalRequest ライバル追加リクエスト
type AddRivalRequest struct {
	Username string `json:"username" validate:"required"`
}

// GetRivals ライバル一覧を取得
func (ctrl *rivalController) GetRivals(c echo.Context) error {
	user := c.Get("user").(*models.User)

	rivals, err := ctrl.rivalUsecase.GetRivals(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ライバル一覧の取得に失敗しました",
		})
	}

	// レスポンス用に変換
	response := make([]map[string]interface{}, len(rivals))
	for i, rival := range rivals {
		response[i] = map[string]interface{}{
			"id":              rival.ID,
			"github_user_id":  rival.RivalGithubUserID,
			"github_username": rival.RivalGithubUsername,
			"avatar_url":      rival.RivalAvatarURL,
			"created_at":      rival.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"rivals":     response,
		"count":      len(rivals),
		"max_rivals": models.MaxRivalsForFreePlan,
	})
}

// AddRival ライバルを追加
func (ctrl *rivalController) AddRival(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var req AddRivalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストが不正です",
		})
	}

	if req.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ユーザー名を指定してください",
		})
	}

	rival, err := ctrl.rivalUsecase.AddRival(c.Request().Context(), user.ID, req.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":              rival.ID,
		"github_user_id":  rival.RivalGithubUserID,
		"github_username": rival.RivalGithubUsername,
		"avatar_url":      rival.RivalAvatarURL,
		"created_at":      rival.CreatedAt,
	})
}

// RemoveRival ライバルを削除
func (ctrl *rivalController) RemoveRival(c echo.Context) error {
	user := c.Get("user").(*models.User)

	rivalIDStr := c.Param("id")
	rivalID, err := strconv.ParseUint(rivalIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "ライバルIDが不正です",
		})
	}

	if err := ctrl.rivalUsecase.RemoveRival(c.Request().Context(), user.ID, rivalID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
