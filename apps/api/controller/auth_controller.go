package controller

import (
	"net/http"

	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
)

// IAuthController 認証コントローラーのインターフェース
type IAuthController interface {
	Callback(c echo.Context) error
	Logout(c echo.Context) error
}

type authController struct {
	userUsecase usecase.IUserUsecase
}

// NewAuthController コンストラクタ
func NewAuthController(userUsecase usecase.IUserUsecase) IAuthController {
	return &authController{
		userUsecase: userUsecase,
	}
}

// CallbackRequest Github OAuth コールバックリクエスト
type CallbackRequest struct {
	GithubUserID   uint64 `json:"github_user_id" validate:"required"`
	GithubUsername string `json:"github_username" validate:"required"`
	Email          string `json:"email"`
	AvatarURL      string `json:"avatar_url"`
}

// Callback Github OAuth コールバック処理
// フロントエンド（NextAuth.js）からユーザー情報を受け取りDBに保存
func (ctrl *authController) Callback(c echo.Context) error {
	var req CallbackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストが不正です",
		})
	}

	if req.GithubUserID == 0 || req.GithubUsername == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Githubユーザー情報が不正です",
		})
	}

	user, err := ctrl.userUsecase.GetOrCreateUser(
		c.Request().Context(),
		req.GithubUserID,
		req.GithubUsername,
		req.Email,
		req.AvatarURL,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "ユーザー情報の保存に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":              user.ID,
		"github_user_id":  user.GithubUserID,
		"github_username": user.GithubUsername,
		"avatar_url":      user.AvatarURL,
	})
}

// Logout ログアウト処理
func (ctrl *authController) Logout(c echo.Context) error {
	// セッションはフロントエンド（NextAuth.js）で管理されるため、
	// バックエンドでは特に処理不要
	return c.JSON(http.StatusOK, map[string]string{
		"message": "ログアウトしました",
	})
}
