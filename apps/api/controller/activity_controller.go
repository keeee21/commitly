package controller

import (
	"net/http"

	"github.com/keeee21/commitly/api/dto"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
)

// IActivityController アクティビティコントローラーのインターフェース
type IActivityController interface {
	GetActivityStream(c echo.Context) error
	GetRhythm(c echo.Context) error
}

type activityController struct {
	activityUsecase usecase.IActivityUsecase
	rivalUsecase    usecase.IRivalUsecase
}

// NewActivityController コンストラクタ
func NewActivityController(activityUsecase usecase.IActivityUsecase, rivalUsecase usecase.IRivalUsecase) IActivityController {
	return &activityController{
		activityUsecase: activityUsecase,
		rivalUsecase:    rivalUsecase,
	}
}

// GetActivityStream アクティビティストリームを取得
// @Summary      アクティビティストリームを取得
// @Description  自分とライバルの直近7日間のコミット活動をストリーム形式で返す
// @Tags         activity
// @Produce      json
// @Success      200 {object} dto.ActivityStreamResponse
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/activity/stream [get]
func (ctrl *activityController) GetActivityStream(c echo.Context) error {
	user := c.Get("user").(*models.User)

	rivals, err := ctrl.rivalUsecase.GetRivals(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "ライバル情報の取得に失敗しました",
		})
	}

	data, err := ctrl.activityUsecase.GetActivityStream(c.Request().Context(), user, rivals)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "アクティビティデータの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, data)
}

// GetRhythm リズム可視化データを取得
// @Summary      リズム可視化データを取得
// @Description  自分とライバルの直近7日間の曜日別コミットパターンを返す
// @Tags         activity
// @Produce      json
// @Success      200 {object} dto.RhythmResponse
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/activity/rhythm [get]
func (ctrl *activityController) GetRhythm(c echo.Context) error {
	user := c.Get("user").(*models.User)

	rivals, err := ctrl.rivalUsecase.GetRivals(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "ライバル情報の取得に失敗しました",
		})
	}

	data, err := ctrl.activityUsecase.GetRhythm(c.Request().Context(), user, rivals)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "リズムデータの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, data)
}
