package router

import (
	"github.com/keeee21/commitly/api/controller"
	"github.com/keeee21/commitly/api/gateway"
	"github.com/keeee21/commitly/api/middleware"
	"github.com/keeee21/commitly/api/repository"
	"github.com/keeee21/commitly/api/usecase"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// SetupRoutes ルーティングをセットアップ
func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	rivalRepo := repository.NewRivalRepository(db)
	commitStatsRepo := repository.NewCommitStatsRepository(db)
	circleRepo := repository.NewCircleRepository(db)
	slackNotificationRepo := repository.NewSlackNotificationSettingRepository(db)

	// Gateways
	githubGateway := gateway.NewGithubGateway("")

	// Usecases
	userUsecase := usecase.NewUserUsecase(userRepo, githubGateway)
	rivalUsecase := usecase.NewRivalUsecase(rivalRepo, githubGateway)
	dashboardUsecase := usecase.NewDashboardUsecase(commitStatsRepo)
	activityUsecase := usecase.NewActivityUsecase(commitStatsRepo)
	circleUsecase := usecase.NewCircleUsecase(circleRepo)
	slackNotificationUsecase := usecase.NewSlackNotificationUsecase(slackNotificationRepo)

	// Controllers
	healthCtrl := controller.NewHealthController()
	authCtrl := controller.NewAuthController(userUsecase)
	userCtrl := controller.NewUserController(userUsecase)
	rivalCtrl := controller.NewRivalController(rivalUsecase)
	dashboardCtrl := controller.NewDashboardController(dashboardUsecase, rivalUsecase)
	activityCtrl := controller.NewActivityController(activityUsecase, rivalUsecase)
	circleCtrl := controller.NewCircleController(circleUsecase)
	slackNotificationCtrl := controller.NewSlackNotificationController(slackNotificationUsecase)

	// Health check
	e.GET("/health", healthCtrl.HealthCheck)

	// API routes
	api := e.Group("/api")

	// Auth routes (認証不要)
	auth := api.Group("/auth")
	auth.POST("/callback", authCtrl.Callback)
	auth.POST("/logout", authCtrl.Logout)

	// Protected routes (認証必要)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(userUsecase))

	// User routes
	protected.GET("/me", userCtrl.GetMe)

	// Rival routes
	rivals := protected.Group("/rivals")
	rivals.GET("", rivalCtrl.GetRivals)
	rivals.POST("", rivalCtrl.AddRival)
	rivals.DELETE("/:id", rivalCtrl.RemoveRival)

	// Dashboard routes
	dashboard := protected.Group("/dashboard")
	dashboard.GET("/weekly", dashboardCtrl.GetWeeklyDashboard)
	dashboard.GET("/monthly", dashboardCtrl.GetMonthlyDashboard)

	// Activity routes
	activity := protected.Group("/activity")
	activity.GET("/stream", activityCtrl.GetActivityStream)
	activity.GET("/rhythm", activityCtrl.GetRhythm)

	// Circle routes
	circles := protected.Group("/circles")
	circles.GET("", circleCtrl.GetCircles)
	circles.POST("", circleCtrl.CreateCircle)
	circles.POST("/join", circleCtrl.JoinCircle)
	circles.DELETE("/:id/leave", circleCtrl.LeaveCircle)
	circles.DELETE("/:id", circleCtrl.DeleteCircle)

	// Slack notification routes
	slack := protected.Group("/notifications/slack")
	slack.GET("", slackNotificationCtrl.GetSetting)
	slack.POST("", slackNotificationCtrl.Create)
	slack.PUT("", slackNotificationCtrl.UpdateEnabled)
	slack.DELETE("", slackNotificationCtrl.Delete)
}
