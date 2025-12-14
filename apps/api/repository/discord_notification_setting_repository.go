package repository

import (
	"context"

	"github.com/keeee21/commitly/api/models"
	"gorm.io/gorm"
)

// IDiscordNotificationSettingRepository Discord通知設定リポジトリのインターフェース
type IDiscordNotificationSettingRepository interface {
	FindByUserID(ctx context.Context, userID uint64) (*models.DiscordNotificationSetting, error)
	FindAllEnabled(ctx context.Context) ([]models.DiscordNotificationSetting, error)
	Upsert(ctx context.Context, setting *models.DiscordNotificationSetting) error
	UpdateEnabled(ctx context.Context, userID uint64, isEnabled bool) error
	Delete(ctx context.Context, userID uint64) error
}

type discordNotificationSettingRepository struct {
	db *gorm.DB
}

// NewDiscordNotificationSettingRepository コンストラクタ
func NewDiscordNotificationSettingRepository(db *gorm.DB) IDiscordNotificationSettingRepository {
	return &discordNotificationSettingRepository{db: db}
}

func (r *discordNotificationSettingRepository) FindByUserID(ctx context.Context, userID uint64) (*models.DiscordNotificationSetting, error) {
	var setting models.DiscordNotificationSetting
	err := r.db.WithContext(ctx).Preload("User").Where("user_id = ?", userID).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

func (r *discordNotificationSettingRepository) FindAllEnabled(ctx context.Context) ([]models.DiscordNotificationSetting, error) {
	var settings []models.DiscordNotificationSetting
	err := r.db.WithContext(ctx).Preload("User").Where("is_enabled = ?", true).Find(&settings).Error
	return settings, err
}

func (r *discordNotificationSettingRepository) Upsert(ctx context.Context, setting *models.DiscordNotificationSetting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

func (r *discordNotificationSettingRepository) UpdateEnabled(ctx context.Context, userID uint64, isEnabled bool) error {
	return r.db.WithContext(ctx).
		Model(&models.DiscordNotificationSetting{}).
		Where("user_id = ?", userID).
		Update("is_enabled", isEnabled).Error
}

func (r *discordNotificationSettingRepository) Delete(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.DiscordNotificationSetting{}).Error
}
