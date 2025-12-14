package repository

import (
	"context"

	"github.com/keeee21/commitly/api/models"
	"gorm.io/gorm"
)

// IRivalRepository ライバルリポジトリのインターフェース
type IRivalRepository interface {
	FindByUserID(ctx context.Context, userID uint64) ([]models.Rival, error)
	FindByID(ctx context.Context, id uint64) (*models.Rival, error)
	FindAllDistinctRivals(ctx context.Context) ([]models.Rival, error)
	CountByUserID(ctx context.Context, userID uint64) (int64, error)
	Create(ctx context.Context, rival *models.Rival) error
	Delete(ctx context.Context, id uint64) error
	ExistsByUserIDAndRivalGithubUserID(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error)
}

type rivalRepository struct {
	db *gorm.DB
}

// NewRivalRepository コンストラクタ
func NewRivalRepository(db *gorm.DB) IRivalRepository {
	return &rivalRepository{db: db}
}

func (r *rivalRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.Rival, error) {
	var rivals []models.Rival
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&rivals).Error; err != nil {
		return nil, err
	}
	return rivals, nil
}

func (r *rivalRepository) FindByID(ctx context.Context, id uint64) (*models.Rival, error) {
	var rival models.Rival
	if err := r.db.WithContext(ctx).First(&rival, id).Error; err != nil {
		return nil, err
	}
	return &rival, nil
}

func (r *rivalRepository) CountByUserID(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Rival{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *rivalRepository) Create(ctx context.Context, rival *models.Rival) error {
	return r.db.WithContext(ctx).Create(rival).Error
}

func (r *rivalRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&models.Rival{}, id).Error
}

func (r *rivalRepository) ExistsByUserIDAndRivalGithubUserID(ctx context.Context, userID uint64, rivalGithubUserID uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Rival{}).
		Where("user_id = ? AND rival_github_user_id = ?", userID, rivalGithubUserID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *rivalRepository) FindAllDistinctRivals(ctx context.Context) ([]models.Rival, error) {
	var rivals []models.Rival
	if err := r.db.WithContext(ctx).
		Distinct("rival_github_user_id", "rival_github_username").
		Find(&rivals).Error; err != nil {
		return nil, err
	}
	return rivals, nil
}
