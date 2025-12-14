package batch

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDateRange_ValidDates(t *testing.T) {
	fromDate, toDate, err := ParseDateRange("2025-01-01", "2025-01-31")

	assert.NoError(t, err)
	assert.NotNil(t, fromDate)
	assert.NotNil(t, toDate)
	assert.Equal(t, 2025, fromDate.Year())
	assert.Equal(t, time.January, fromDate.Month())
	assert.Equal(t, 1, fromDate.Day())
	assert.Equal(t, 2025, toDate.Year())
	assert.Equal(t, time.January, toDate.Month())
	assert.Equal(t, 31, toDate.Day())
}

func TestParseDateRange_OnlyFromDate(t *testing.T) {
	fromDate, toDate, err := ParseDateRange("2025-01-01", "")

	assert.NoError(t, err)
	assert.NotNil(t, fromDate)
	assert.Nil(t, toDate)
	assert.Equal(t, 2025, fromDate.Year())
}

func TestParseDateRange_OnlyToDate(t *testing.T) {
	fromDate, toDate, err := ParseDateRange("", "2025-01-31")

	assert.NoError(t, err)
	assert.Nil(t, fromDate)
	assert.NotNil(t, toDate)
	assert.Equal(t, 2025, toDate.Year())
}

func TestParseDateRange_EmptyDates(t *testing.T) {
	fromDate, toDate, err := ParseDateRange("", "")

	assert.NoError(t, err)
	assert.Nil(t, fromDate)
	assert.Nil(t, toDate)
}

func TestParseDateRange_InvalidFromDate(t *testing.T) {
	fromDate, toDate, err := ParseDateRange("invalid-date", "2025-01-31")

	assert.Error(t, err)
	assert.Nil(t, fromDate)
	assert.Nil(t, toDate)
	assert.Contains(t, err.Error(), "invalid from date format")
}

func TestParseDateRange_InvalidToDate(t *testing.T) {
	fromDate, toDate, err := ParseDateRange("2025-01-01", "invalid-date")

	assert.Error(t, err)
	assert.Nil(t, fromDate)
	assert.Nil(t, toDate)
	assert.Contains(t, err.Error(), "invalid to date format")
}

func TestParseDateRange_WrongFormat(t *testing.T) {
	// 日付フォーマットが異なる場合
	fromDate, toDate, err := ParseDateRange("01/01/2025", "")

	assert.Error(t, err)
	assert.Nil(t, fromDate)
	assert.Nil(t, toDate)
}

func TestValidateDateRange_ValidRange(t *testing.T) {
	fromDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	err := ValidateDateRange(&fromDate, &toDate)

	assert.NoError(t, err)
}

func TestValidateDateRange_SameDate(t *testing.T) {
	date := time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC)

	err := ValidateDateRange(&date, &date)

	assert.NoError(t, err)
}

func TestValidateDateRange_FromAfterTo(t *testing.T) {
	fromDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	err := ValidateDateRange(&fromDate, &toDate)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "from date must be before to date")
}

func TestValidateDateRange_NilDates(t *testing.T) {
	err := ValidateDateRange(nil, nil)

	assert.NoError(t, err)
}

func TestValidateDateRange_NilFromDate(t *testing.T) {
	toDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	err := ValidateDateRange(nil, &toDate)

	assert.NoError(t, err)
}

func TestValidateDateRange_NilToDate(t *testing.T) {
	fromDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	err := ValidateDateRange(&fromDate, nil)

	assert.NoError(t, err)
}

// mockSyncCommitsUsecase テスト用のモック
type mockSyncCommitsUsecase struct {
	SyncAllUsersFunc              func(ctx context.Context) error
	SyncAllUsersWithDateRangeFunc func(ctx context.Context, fromDate, toDate *time.Time) error
	SyncUserFunc                  func(ctx context.Context, githubUserID uint64, githubUsername string, fromDate, toDate *time.Time) error
}

func (m *mockSyncCommitsUsecase) SyncAllUsers(ctx context.Context) error {
	if m.SyncAllUsersFunc != nil {
		return m.SyncAllUsersFunc(ctx)
	}
	return nil
}

func (m *mockSyncCommitsUsecase) SyncAllUsersWithDateRange(ctx context.Context, fromDate, toDate *time.Time) error {
	if m.SyncAllUsersWithDateRangeFunc != nil {
		return m.SyncAllUsersWithDateRangeFunc(ctx, fromDate, toDate)
	}
	return nil
}

func (m *mockSyncCommitsUsecase) SyncUser(ctx context.Context, githubUserID uint64, githubUsername string, fromDate, toDate *time.Time) error {
	if m.SyncUserFunc != nil {
		return m.SyncUserFunc(ctx, githubUserID, githubUsername, fromDate, toDate)
	}
	return nil
}

func TestRunSyncCommits_Success(t *testing.T) {
	ctx := context.Background()
	called := false

	mockUsecase := &mockSyncCommitsUsecase{
		SyncAllUsersWithDateRangeFunc: func(ctx context.Context, fromDate, toDate *time.Time) error {
			called = true
			return nil
		},
	}

	config := SyncCommitsConfig{
		FromDate: "2025-01-01",
		ToDate:   "2025-01-31",
	}

	err := RunSyncCommits(ctx, mockUsecase, config)

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestRunSyncCommits_WithoutDates(t *testing.T) {
	ctx := context.Background()
	var capturedFrom, capturedTo *time.Time

	mockUsecase := &mockSyncCommitsUsecase{
		SyncAllUsersWithDateRangeFunc: func(ctx context.Context, fromDate, toDate *time.Time) error {
			capturedFrom = fromDate
			capturedTo = toDate
			return nil
		},
	}

	config := SyncCommitsConfig{}

	err := RunSyncCommits(ctx, mockUsecase, config)

	assert.NoError(t, err)
	assert.Nil(t, capturedFrom)
	assert.Nil(t, capturedTo)
}

func TestRunSyncCommits_InvalidFromDate(t *testing.T) {
	ctx := context.Background()

	mockUsecase := &mockSyncCommitsUsecase{}

	config := SyncCommitsConfig{
		FromDate: "invalid",
	}

	err := RunSyncCommits(ctx, mockUsecase, config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid from date format")
}

func TestRunSyncCommits_InvalidToDate(t *testing.T) {
	ctx := context.Background()

	mockUsecase := &mockSyncCommitsUsecase{}

	config := SyncCommitsConfig{
		ToDate: "invalid",
	}

	err := RunSyncCommits(ctx, mockUsecase, config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid to date format")
}

func TestRunSyncCommits_FromAfterTo(t *testing.T) {
	ctx := context.Background()

	mockUsecase := &mockSyncCommitsUsecase{}

	config := SyncCommitsConfig{
		FromDate: "2025-02-01",
		ToDate:   "2025-01-01",
	}

	err := RunSyncCommits(ctx, mockUsecase, config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "from date must be before to date")
}

func TestRunSyncCommits_UsecaseError(t *testing.T) {
	ctx := context.Background()

	mockUsecase := &mockSyncCommitsUsecase{
		SyncAllUsersWithDateRangeFunc: func(ctx context.Context, fromDate, toDate *time.Time) error {
			return errors.New("sync failed")
		},
	}

	config := SyncCommitsConfig{}

	err := RunSyncCommits(ctx, mockUsecase, config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to sync commits")
}

func TestRunSyncCommits_PassesDatesToUsecase(t *testing.T) {
	ctx := context.Background()
	var capturedFrom, capturedTo *time.Time

	mockUsecase := &mockSyncCommitsUsecase{
		SyncAllUsersWithDateRangeFunc: func(ctx context.Context, fromDate, toDate *time.Time) error {
			capturedFrom = fromDate
			capturedTo = toDate
			return nil
		},
	}

	config := SyncCommitsConfig{
		FromDate: "2025-06-01",
		ToDate:   "2025-06-30",
	}

	err := RunSyncCommits(ctx, mockUsecase, config)

	assert.NoError(t, err)
	assert.NotNil(t, capturedFrom)
	assert.NotNil(t, capturedTo)
	assert.Equal(t, 2025, capturedFrom.Year())
	assert.Equal(t, time.June, capturedFrom.Month())
	assert.Equal(t, 1, capturedFrom.Day())
	assert.Equal(t, 2025, capturedTo.Year())
	assert.Equal(t, time.June, capturedTo.Month())
	assert.Equal(t, 30, capturedTo.Day())
}
