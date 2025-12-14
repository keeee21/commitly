package middleware

import (
	"net/http"
	"strconv"

	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware 認証ミドルウェア
func AuthMiddleware(userUsecase usecase.IUserUsecase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// X-GitHub-User-ID ヘッダーからGithubユーザーIDを取得
			// フロントエンド（NextAuth.js）のセッションから取得してヘッダーに設定される想定
			githubUserIDStr := c.Request().Header.Get("X-GitHub-User-ID")
			if githubUserIDStr == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "認証が必要です",
				})
			}

			githubUserID, err := strconv.ParseUint(githubUserIDStr, 10, 64)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "認証情報が不正です",
				})
			}

			// ユーザーをDBから取得
			user, err := userUsecase.GetUserByGithubUserID(c.Request().Context(), githubUserID)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "ユーザーが見つかりません",
				})
			}

			// コンテキストにユーザー情報を設定
			c.Set("user", user)

			return next(c)
		}
	}
}
