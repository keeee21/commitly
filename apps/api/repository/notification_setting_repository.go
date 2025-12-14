package repository

import (
	"context"

	"github.com/keeee21/commitly/api/models"
	"gorm.io/gorm"
)

// INotificationSettingRepository 通知設定リポジトリのインターフェース
type INotificationSettingRepository interface {
	FindByUserID(ctx context.Context, userID uint64) ([]models.NotificationSetting, error)
	FindByID(ctx context.Context, id uint64) (*models.NotificationSetting, error)
	FindEnabledSettings(ctx context.Context) ([]models.NotificationSetting, error)
	Create(ctx context.Context, setting *models.NotificationSetting) error
	Update(ctx context.Context, setting *models.NotificationSetting) error
	Delete(ctx context.Context, id uint64) error
}

type notificationSettingRepository struct {
	db *gorm.DB
}

// NewNotificationSettingRepository コンストラクタ
func NewNotificationSettingRepository(db *gorm.DB) INotificationSettingRepository {
	return &notificationSettingRepository{db: db}
}

func (r *notificationSettingRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
	var settings []models.NotificationSetting
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *notificationSettingRepository) FindByID(ctx context.Context, id uint64) (*models.NotificationSetting, error) {
	var setting models.NotificationSetting
	if err := r.db.WithContext(ctx).First(&setting, id).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *notificationSettingRepository) FindEnabledSettings(ctx context.Context) ([]models.NotificationSetting, error) {
	var settings []models.NotificationSetting
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("is_enabled = ?", true).
		Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *notificationSettingRepository) Create(ctx context.Context, setting *models.NotificationSetting) error {
	return r.db.WithContext(ctx).Create(setting).Error
}

func (r *notificationSettingRepository) Update(ctx context.Context, setting *models.NotificationSetting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

func (r *notificationSettingRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&models.NotificationSetting{}, id).Error
}
