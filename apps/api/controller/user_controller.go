package controller

import (
	"net/http"

	"github.com/keeee21/commitly/api/dto"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
)

// IUserController ユーザーコントローラーのインターフェース
type IUserController interface {
	GetMe(c echo.Context) error
}

type userController struct {
	userUsecase usecase.IUserUsecase
}

// NewUserController コンストラクタ
func NewUserController(userUsecase usecase.IUserUsecase) IUserController {
	return &userController{
		userUsecase: userUsecase,
	}
}

// GetMe 現在のユーザー情報を取得
// @Summary      現在のユーザー情報を取得
// @Description  認証済みユーザーの情報を返す
// @Tags         user
// @Produce      json
// @Success      200 {object} dto.UserResponse
// @Security     GitHubUserID
// @Router       /api/me [get]
func (ctrl *userController) GetMe(c echo.Context) error {
	user := c.Get("user").(*models.User)

	return c.JSON(http.StatusOK, dto.UserResponse{
		ID:             user.ID,
		GithubUserID:   user.GithubUserID,
		GithubUsername: user.GithubUsername,
		AvatarURL:      user.AvatarURL,
		CreatedAt:      user.CreatedAt,
	})
}
