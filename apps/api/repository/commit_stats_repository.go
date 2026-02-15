package repository

import (
	"context"
	"time"

	"github.com/keeee21/commitly/api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ICommitStatsRepository コミット統計リポジトリのインターフェース
type ICommitStatsRepository interface {
	// 日別のコミット統計を取得（期間指定）
	FindByGithubUserIDAndDateRange(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error)
	// 複数ユーザーの日別コミット統計を取得
	FindByGithubUserIDsAndDateRange(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error)
	// コミット統計を保存（Upsert）
	Upsert(ctx context.Context, stats *models.CommitStats) error
	// バッチでコミット統計を保存
	UpsertBatch(ctx context.Context, statsList []models.CommitStats) error
}

type commitStatsRepository struct {
	db *gorm.DB
}

// NewCommitStatsRepository コンストラクタ
func NewCommitStatsRepository(db *gorm.DB) ICommitStatsRepository {
	return &commitStatsRepository{db: db}
}

func (r *commitStatsRepository) FindByGithubUserIDAndDateRange(ctx context.Context, githubUserID uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	var stats []models.CommitStats
	if err := r.db.WithContext(ctx).
		Where("github_user_id = ? AND date >= ? AND date <= ?", githubUserID, startDate, endDate).
		Order("date ASC, repository ASC").
		Find(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *commitStatsRepository) FindByGithubUserIDsAndDateRange(ctx context.Context, githubUserIDs []uint64, startDate, endDate time.Time) ([]models.CommitStats, error) {
	var stats []models.CommitStats
	if err := r.db.WithContext(ctx).
		Where("github_user_id IN ? AND date >= ? AND date <= ?", githubUserIDs, startDate, endDate).
		Order("github_user_id ASC, date ASC, repository ASC").
		Find(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *commitStatsRepository) Upsert(ctx context.Context, stats *models.CommitStats) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "github_user_id"}, {Name: "date"}, {Name: "repository"}},
		DoUpdates: clause.AssignmentColumns([]string{"commit_count", "primary_hour", "language", "fetched_at"}),
	}).Create(stats).Error
}

func (r *commitStatsRepository) UpsertBatch(ctx context.Context, statsList []models.CommitStats) error {
	if len(statsList) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "github_user_id"}, {Name: "date"}, {Name: "repository"}},
		DoUpdates: clause.AssignmentColumns([]string{"commit_count", "primary_hour", "language", "fetched_at"}),
	}).Create(&statsList).Error
}
