package controller

import (
	"net/http"
	"strconv"

	"github.com/keeee21/commitly/api/dto"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
)

// ICircleController サークルコントローラーのインターフェース
type ICircleController interface {
	GetCircles(c echo.Context) error
	CreateCircle(c echo.Context) error
	JoinCircle(c echo.Context) error
	LeaveCircle(c echo.Context) error
	DeleteCircle(c echo.Context) error
}

type circleController struct {
	circleUsecase usecase.ICircleUsecase
}

// NewCircleController コンストラクタ
func NewCircleController(circleUsecase usecase.ICircleUsecase) ICircleController {
	return &circleController{
		circleUsecase: circleUsecase,
	}
}

// GetCircles サークル一覧を取得
// @Summary      サークル一覧を取得
// @Description  自分が所属するサークル一覧を返す
// @Tags         circles
// @Produce      json
// @Success      200 {object} dto.CirclesListResponse
// @Failure      500 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/circles [get]
func (ctrl *circleController) GetCircles(c echo.Context) error {
	user := c.Get("user").(*models.User)

	circles, err := ctrl.circleUsecase.GetCircles(c.Request().Context(), user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "サークル一覧の取得に失敗しました",
		})
	}

	response := make([]dto.CircleResponse, len(circles))
	for i, circle := range circles {
		response[i] = buildCircleResponse(circle, user.ID)
	}

	return c.JSON(http.StatusOK, dto.CirclesListResponse{
		Circles:    response,
		Count:      len(circles),
		MaxCircles: models.MaxCirclesForFreePlan,
	})
}

// CreateCircle サークルを作成
// @Summary      サークルを作成
// @Description  新しいサークルを作成する
// @Tags         circles
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateCircleRequest true "サークル作成リクエスト"
// @Success      201 {object} dto.CircleResponse
// @Failure      400 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/circles [post]
func (ctrl *circleController) CreateCircle(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var req dto.CreateCircleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "リクエストが不正です",
		})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "サークル名を入力してください",
		})
	}

	circle, err := ctrl.circleUsecase.CreateCircle(c.Request().Context(), user.ID, req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, buildCircleResponse(*circle, user.ID))
}

// JoinCircle サークルに参加
// @Summary      サークルに参加
// @Description  招待コードを使ってサークルに参加する
// @Tags         circles
// @Accept       json
// @Produce      json
// @Param        request body dto.JoinCircleRequest true "サークル参加リクエスト"
// @Success      200 {object} dto.CircleResponse
// @Failure      400 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/circles/join [post]
func (ctrl *circleController) JoinCircle(c echo.Context) error {
	user := c.Get("user").(*models.User)

	var req dto.JoinCircleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "リクエストが不正です",
		})
	}

	if req.InviteCode == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "招待コードを入力してください",
		})
	}

	circle, err := ctrl.circleUsecase.JoinCircle(c.Request().Context(), user.ID, req.InviteCode)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, buildCircleResponse(*circle, user.ID))
}

// LeaveCircle サークルを退会
// @Summary      サークルを退会
// @Description  指定IDのサークルから退会する
// @Tags         circles
// @Param        id path int true "サークルID"
// @Success      204
// @Failure      400 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/circles/{id}/leave [delete]
func (ctrl *circleController) LeaveCircle(c echo.Context) error {
	user := c.Get("user").(*models.User)

	circleIDStr := c.Param("id")
	circleID, err := strconv.ParseUint(circleIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "サークルIDが不正です",
		})
	}

	if err := ctrl.circleUsecase.LeaveCircle(c.Request().Context(), user.ID, circleID); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteCircle サークルを削除
// @Summary      サークルを削除
// @Description  指定IDのサークルを削除する（オーナーのみ）
// @Tags         circles
// @Param        id path int true "サークルID"
// @Success      204
// @Failure      400 {object} dto.ErrorResponse
// @Security     GitHubUserID
// @Router       /api/circles/{id} [delete]
func (ctrl *circleController) DeleteCircle(c echo.Context) error {
	user := c.Get("user").(*models.User)

	circleIDStr := c.Param("id")
	circleID, err := strconv.ParseUint(circleIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "サークルIDが不正です",
		})
	}

	if err := ctrl.circleUsecase.DeleteCircle(c.Request().Context(), user.ID, circleID); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func buildCircleResponse(circle models.Circle, currentUserID uint64) dto.CircleResponse {
	members := make([]dto.CircleMemberResponse, len(circle.Members))
	for j, member := range circle.Members {
		members[j] = dto.CircleMemberResponse{
			GithubUsername: member.User.GithubUsername,
			AvatarURL:      member.User.AvatarURL,
			JoinedAt:       member.JoinedAt,
		}
	}

	return dto.CircleResponse{
		ID:         circle.ID,
		Name:       circle.Name,
		InviteCode: circle.InviteCode,
		IsOwner:    circle.OwnerUserID == currentUserID,
		Members:    members,
		CreatedAt:  circle.CreatedAt,
	}
}
