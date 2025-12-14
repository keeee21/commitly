package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/keeee21/commitly/api/batch"
	"github.com/keeee21/commitly/api/db"
	"github.com/keeee21/commitly/api/gateway"
	"github.com/keeee21/commitly/api/repository"
	"github.com/keeee21/commitly/api/usecase"
)

func main() {
	// Parse command line flags
	command := flag.String("command", "", "batch command to run (sync-commits, send-notifications)")
	fromDate := flag.String("from", "", "start date for sync (YYYY-MM-DD)")
	toDate := flag.String("to", "", "end date for sync (YYYY-MM-DD)")
	period := flag.String("period", "weekly", "notification period (weekly, monthly)")
	flag.Parse()

	if *command == "" {
		log.Fatal("command flag is required. Available commands: sync-commits, send-notifications")
	}

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	database, err := db.NewDatabase(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()

	// Run command
	switch *command {
	case "sync-commits":
		// Initialize repositories
		userRepo := repository.NewUserRepository(database)
		rivalRepo := repository.NewRivalRepository(database)
		commitStatsRepo := repository.NewCommitStatsRepository(database)

		// Initialize gateway with GitHub token
		githubToken := os.Getenv("GITHUB_TOKEN")
		githubGateway := gateway.NewGithubGateway(githubToken)

		// Initialize usecase
		syncUsecase := usecase.NewSyncCommitsUsecase(userRepo, rivalRepo, commitStatsRepo, githubGateway)

		// Run sync
		config := batch.SyncCommitsConfig{
			FromDate: *fromDate,
			ToDate:   *toDate,
		}
		if err := batch.RunSyncCommits(ctx, syncUsecase, config); err != nil {
			log.Fatalf("Failed to run sync-commits: %v", err)
		}

	case "send-notifications":
		// Initialize repositories
		slackNotificationRepo := repository.NewSlackNotificationSettingRepository(database)
		notificationLogRepo := repository.NewNotificationLogRepository(database)
		rivalRepo := repository.NewRivalRepository(database)
		commitStatsRepo := repository.NewCommitStatsRepository(database)

		// Initialize gateway
		slackGateway := gateway.NewSlackGateway()

		// Initialize dependencies
		deps := &batch.SendNotificationsDeps{
			SlackNotificationRepo: slackNotificationRepo,
			NotificationLogRepo:   notificationLogRepo,
			RivalRepo:             rivalRepo,
			CommitStatsRepo:       commitStatsRepo,
			SlackGateway:          slackGateway,
		}

		// Run send notifications
		config := batch.SendNotificationsConfig{
			Period: *period,
		}
		if err := batch.RunSendNotifications(ctx, deps, config); err != nil {
			log.Fatalf("Failed to run send-notifications: %v", err)
		}

	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}
