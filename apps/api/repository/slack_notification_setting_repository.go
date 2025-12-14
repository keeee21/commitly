package repository

import (
	"context"

	"github.com/keeee21/commitly/api/models"
	"gorm.io/gorm"
)

// ISlackNotificationSettingRepository Slack通知設定リポジトリのインターフェース
type ISlackNotificationSettingRepository interface {
	FindByUserID(ctx context.Context, userID uint64) (*models.SlackNotificationSetting, error)
	FindAllEnabled(ctx context.Context) ([]models.SlackNotificationSetting, error)
	Upsert(ctx context.Context, setting *models.SlackNotificationSetting) error
	UpdateEnabled(ctx context.Context, userID uint64, isEnabled bool) error
	Delete(ctx context.Context, userID uint64) error
}

type slackNotificationSettingRepository struct {
	db *gorm.DB
}

// NewSlackNotificationSettingRepository コンストラクタ
func NewSlackNotificationSettingRepository(db *gorm.DB) ISlackNotificationSettingRepository {
	return &slackNotificationSettingRepository{db: db}
}

func (r *slackNotificationSettingRepository) FindByUserID(ctx context.Context, userID uint64) (*models.SlackNotificationSetting, error) {
	var setting models.SlackNotificationSetting
	err := r.db.WithContext(ctx).Preload("User").Where("user_id = ?", userID).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

func (r *slackNotificationSettingRepository) FindAllEnabled(ctx context.Context) ([]models.SlackNotificationSetting, error) {
	var settings []models.SlackNotificationSetting
	err := r.db.WithContext(ctx).Preload("User").Where("is_enabled = ?", true).Find(&settings).Error
	return settings, err
}

func (r *slackNotificationSettingRepository) Upsert(ctx context.Context, setting *models.SlackNotificationSetting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

func (r *slackNotificationSettingRepository) UpdateEnabled(ctx context.Context, userID uint64, isEnabled bool) error {
	return r.db.WithContext(ctx).
		Model(&models.SlackNotificationSetting{}).
		Where("user_id = ?", userID).
		Update("is_enabled", isEnabled).Error
}

func (r *slackNotificationSettingRepository) Delete(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.SlackNotificationSetting{}).Error
}
