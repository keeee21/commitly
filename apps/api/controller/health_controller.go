package controller

import (
	"net/http"

	"github.com/keeee21/commitly/api/dto"
	"github.com/labstack/echo/v4"
)

// IHealthController ヘルスチェックコントローラーのインターフェース
type IHealthController interface {
	HealthCheck(c echo.Context) error
}

type healthController struct{}

// NewHealthController コンストラクタ
func NewHealthController() IHealthController {
	return &healthController{}
}

// HealthCheck ヘルスチェック
// @Summary      ヘルスチェック
// @Description  サーバーの稼働状態を確認する
// @Tags         health
// @Produce      json
// @Success      200 {object} dto.HealthResponse
// @Router       /health [get]
func (ctrl *healthController) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, dto.HealthResponse{
		Status: "ok",
	})
}
