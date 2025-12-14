package controller

import (
	"net/http"

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
func (ctrl *userController) GetMe(c echo.Context) error {
	user := c.Get("user").(*models.User)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":              user.ID,
		"github_user_id":  user.GithubUserID,
		"github_username": user.GithubUsername,
		"avatar_url":      user.AvatarURL,
		"created_at":      user.CreatedAt,
	})
}
