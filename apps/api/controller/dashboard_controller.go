package controller

import (
	"net/http"

	"github.com/keeee21/commitly/api/dto"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
)

// IDashboardController ダッシュボードコントローラーのインターフェース
type IDashboardController interface {
	GetWeeklyDashboard(c echo.Context) error
	GetMonthlyDashboard(c echo.Context) error
}

type dashboardController struct {
	dashboardUsecase usecase.IDashboardUsecase
	rivalUsecase     usecase.IRivalUsecase
}

// NewDashboardController コンストラクタ
func NewDashboardController(dashboardUsecase usecase.IDashboardUsecase, rivalUsecase usecase.IRivalUsecase) IDashboardController {
	return &dashboardController{
		dashboardUsecase: dashboardUsecase,
		rivalUsecase:     rivalUsecase,
	}
}

// GetWeeklyDashboard 週間ダッシュボードデータを取得
// @Summary      週間ダッシュボードデータを取得
// @Description  直近7日間のコミット統計を自分とライバルで比較
// @Tags         dashboard
// @Produce      json
// @Success      200 {object} usecase.DashboardData
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/dashboard/weekly [get]
func (ctrl *dashboardController) GetWeeklyDashboard(c echo.Context) error {
	user := c.Get("user").(*models.User)

	rivals, err := ctrl.rivalUsecase.GetRivals(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "ライバル情報の取得に失敗しました",
		})
	}

	data, err := ctrl.dashboardUsecase.GetWeeklyDashboard(c.Request().Context(), user, rivals)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "ダッシュボードデータの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, data)
}

// GetMonthlyDashboard 月間ダッシュボードデータを取得
// @Summary      月間ダッシュボードデータを取得
// @Description  今月のコミット統計を自分とライバルで比較
// @Tags         dashboard
// @Produce      json
// @Success      200 {object} usecase.DashboardData
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/dashboard/monthly [get]
func (ctrl *dashboardController) GetMonthlyDashboard(c echo.Context) error {
	user := c.Get("user").(*models.User)

	rivals, err := ctrl.rivalUsecase.GetRivals(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "ライバル情報の取得に失敗しました",
		})
	}

	data, err := ctrl.dashboardUsecase.GetMonthlyDashboard(c.Request().Context(), user, rivals)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "ダッシュボードデータの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, data)
}
