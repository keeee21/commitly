package controller

import (
	"net/http"
	"strconv"

	"github.com/keeee21/commitly/api/dto"
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

// GetRivals ライバル一覧を取得
// @Summary      ライバル一覧を取得
// @Description  登録済みのライバル一覧を返す
// @Tags         rivals
// @Produce      json
// @Success      200 {object} dto.RivalsListResponse
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/rivals [get]
func (ctrl *rivalController) GetRivals(c echo.Context) error {
	user := c.Get("user").(*models.User)

	rivals, err := ctrl.rivalUsecase.GetRivals(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "ライバル一覧の取得に失敗しました",
		})
	}

	response := make([]dto.RivalResponse, len(rivals))
	for i, rival := range rivals {
		response[i] = dto.RivalResponse{
			ID:             rival.ID,
			GithubUserID:   rival.RivalGithubUserID,
			GithubUsername: rival.RivalGithubUsername,
			AvatarURL:      rival.RivalAvatarURL,
			CreatedAt:      rival.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, dto.RivalsListResponse{
		Rivals:    response,
		Count:     len(rivals),
		MaxRivals: models.MaxRivalsForFreePlan,
	})
}

// AddRival ライバルを追加
// @Summary      ライバルを追加
// @Description  GitHubユーザー名を指定してライバルを登録する
// @Tags         rivals
// @Accept       json
// @Produce      json
// @Param        request body dto.AddRivalRequest true "ライバル追加リクエスト"
// @Success      201 {object} dto.RivalResponse
// @Failure      400 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/rivals [post]
func (ctrl *rivalController) AddRival(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var req dto.AddRivalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "リクエストが不正です",
		})
	}

	if req.Username == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "ユーザー名を指定してください",
		})
	}

	rival, err := ctrl.rivalUsecase.AddRival(c.Request().Context(), user.ID, req.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.RivalResponse{
		ID:             rival.ID,
		GithubUserID:   rival.RivalGithubUserID,
		GithubUsername: rival.RivalGithubUsername,
		AvatarURL:      rival.RivalAvatarURL,
		CreatedAt:      rival.CreatedAt,
	})
}

// RemoveRival ライバルを削除
// @Summary      ライバルを削除
// @Description  指定IDのライバルを削除する
// @Tags         rivals
// @Param        id path int true "ライバルID"
// @Success      204
// @Failure      400 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/rivals/{id} [delete]
func (ctrl *rivalController) RemoveRival(c echo.Context) error {
	user := c.Get("user").(*models.User)

	rivalIDStr := c.Param("id")
	rivalID, err := strconv.ParseUint(rivalIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "ライバルIDが不正です",
		})
	}

	if err := ctrl.rivalUsecase.RemoveRival(c.Request().Context(), user.ID, rivalID); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
