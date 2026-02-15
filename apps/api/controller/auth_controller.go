package controller

import (
	"net/http"

	"github.com/keeee21/commitly/api/dto"
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

// Callback Github OAuth コールバック処理
// @Summary      Github OAuth コールバック
// @Description  フロントエンド（NextAuth.js）からユーザー情報を受け取りDBに保存
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.CallbackRequest true "コールバックリクエスト"
// @Success      200 {object} dto.UserResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /api/auth/callback [post]
func (ctrl *authController) Callback(c echo.Context) error {
	var req dto.CallbackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "リクエストが不正です",
		})
	}

	if req.GithubUserID == 0 || req.GithubUsername == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Githubユーザー情報が不正です",
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
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "ユーザー情報の保存に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		ID:             user.ID,
		GithubUserID:   user.GithubUserID,
		GithubUsername: user.GithubUsername,
		AvatarURL:      user.AvatarURL,
	})
}

// Logout ログアウト処理
// @Summary      ログアウト
// @Description  セッションはフロントエンド（NextAuth.js）で管理されるため、バックエンドでは特に処理不要
// @Tags         auth
// @Produce      json
// @Success      200 {object} dto.MessageResponse
// @Router       /api/auth/logout [post]
func (ctrl *authController) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, dto.MessageResponse{
		Message: "ログアウトしました",
	})
}
