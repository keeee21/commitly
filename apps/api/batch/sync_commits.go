package batch

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/keeee21/commitly/api/usecase"
)

// SyncCommitsConfig sync-commitsバッチの設定
type SyncCommitsConfig struct {
	FromDate string
	ToDate   string
}

// ParseDateRange 日付範囲をパースする
func ParseDateRange(fromDateStr, toDateStr string) (*time.Time, *time.Time, error) {
	var fromDate, toDate *time.Time

	if fromDateStr != "" {
		t, err := time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid from date format: %w", err)
		}
		fromDate = &t
	}

	if toDateStr != "" {
		t, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid to date format: %w", err)
		}
		toDate = &t
	}

	return fromDate, toDate, nil
}

// ValidateDateRange 日付範囲の妥当性を検証する
func ValidateDateRange(fromDate, toDate *time.Time) error {
	if fromDate != nil && toDate != nil {
		if fromDate.After(*toDate) {
			return fmt.Errorf("from date must be before to date")
		}
	}
	return nil
}

// RunSyncCommits sync-commitsバッチを実行する
func RunSyncCommits(ctx context.Context, syncUsecase usecase.ISyncCommitsUsecase, config SyncCommitsConfig) error {
	log.Println("Starting sync-commits batch...")
	startTime := time.Now()

	// Parse date options
	fromDate, toDate, err := ParseDateRange(config.FromDate, config.ToDate)
	if err != nil {
		return err
	}

	// Validate date range
	if err := ValidateDateRange(fromDate, toDate); err != nil {
		return err
	}

	if fromDate != nil {
		log.Printf("From date: %s", fromDate.Format("2006-01-02"))
	}
	if toDate != nil {
		log.Printf("To date: %s", toDate.Format("2006-01-02"))
	}

	// Run sync
	if err := syncUsecase.SyncAllUsersWithDateRange(ctx, fromDate, toDate); err != nil {
		return fmt.Errorf("failed to sync commits: %w", err)
	}

	elapsed := time.Since(startTime)
	log.Printf("sync-commits batch completed in %s", elapsed)

	return nil
}
