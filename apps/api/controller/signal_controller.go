package controller

import (
	"net/http"
	"strconv"

	"github.com/keeee21/commitly/api/dto"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
)

// ISignalController シグナルコントローラーのインターフェース
type ISignalController interface {
	GetSignals(c echo.Context) error
	GetRecentSignals(c echo.Context) error
}

type signalController struct {
	signalUsecase usecase.ISignalUsecase
}

// NewSignalController コンストラクタ
func NewSignalController(signalUsecase usecase.ISignalUsecase) ISignalController {
	return &signalController{
		signalUsecase: signalUsecase,
	}
}

// GetSignals サークルのシグナルを取得
// @Summary      サークルの並走シグナルを取得
// @Description  指定サークルのメンバー間の並走シグナルを返す
// @Tags         signals
// @Produce      json
// @Param        id path int true "サークルID"
// @Success      200 {object} dto.SignalsListResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/circles/{id}/signals [get]
func (ctrl *signalController) GetSignals(c echo.Context) error {
	user := c.Get("user").(*models.User)

	circleIDStr := c.Param("id")
	circleID, err := strconv.ParseUint(circleIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "サークルIDが不正です",
		})
	}

	signals, err := ctrl.signalUsecase.GetSignals(c.Request().Context(), user.ID, circleID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, buildSignalsListResponse(signals))
}

// GetRecentSignals 全サークル横断の最近のシグナルを取得
// @Summary      最近の並走シグナルを取得
// @Description  全サークル横断で直近の並走シグナルを返す
// @Tags         signals
// @Produce      json
// @Success      200 {object} dto.SignalsListResponse
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/signals/recent [get]
func (ctrl *signalController) GetRecentSignals(c echo.Context) error {
	user := c.Get("user").(*models.User)

	signals, err := ctrl.signalUsecase.GetRecentSignals(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "シグナルの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, buildSignalsListResponse(signals))
}

func buildSignalsListResponse(signals []usecase.Signal) dto.SignalsListResponse {
	response := make([]dto.SignalResponse, len(signals))
	for i, s := range signals {
		users := make([]dto.SignalUserResponse, len(s.Usernames))
		for j := range s.Usernames {
			users[j] = dto.SignalUserResponse{
				GithubUsername: s.Usernames[j],
				AvatarURL:      s.AvatarURLs[j],
			}
		}
		response[i] = dto.SignalResponse{
			Type:       s.Type,
			Date:       s.Date,
			Users:      users,
			Detail:     s.Detail,
			CircleID:   s.CircleID,
			CircleName: s.CircleName,
		}
	}
	return dto.SignalsListResponse{
		Signals: response,
	}
}
