package controller

import (
	"net/http"

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
func (ctrl *dashboardController) GetWeeklyDashboard(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// ライバル一覧を取得
	rivals, err := ctrl.rivalUsecase.GetRivals(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ライバル情報の取得に失敗しました",
		})
	}

	// ダッシュボードデータを取得
	data, err := ctrl.dashboardUsecase.GetWeeklyDashboard(c.Request().Context(), user, rivals)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ダッシュボードデータの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, data)
}

// GetMonthlyDashboard 月間ダッシュボードデータを取得
func (ctrl *dashboardController) GetMonthlyDashboard(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// ライバル一覧を取得
	rivals, err := ctrl.rivalUsecase.GetRivals(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ライバル情報の取得に失敗しました",
		})
	}

	// ダッシュボードデータを取得
	data, err := ctrl.dashboardUsecase.GetMonthlyDashboard(c.Request().Context(), user, rivals)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ダッシュボードデータの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, data)
}
