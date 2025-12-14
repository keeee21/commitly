package repository

import (
	"context"

	"github.com/keeee21/commitly/api/models"
	"gorm.io/gorm"
)

// ILineNotificationSettingRepository LINE通知設定リポジトリのインターフェース
type ILineNotificationSettingRepository interface {
	FindByUserID(ctx context.Context, userID uint64) (*models.LineNotificationSetting, error)
	FindAllEnabled(ctx context.Context) ([]models.LineNotificationSetting, error)
	Upsert(ctx context.Context, setting *models.LineNotificationSetting) error
	UpdateEnabled(ctx context.Context, userID uint64, isEnabled bool) error
	Delete(ctx context.Context, userID uint64) error
}

type lineNotificationSettingRepository struct {
	db *gorm.DB
}

// NewLineNotificationSettingRepository コンストラクタ
func NewLineNotificationSettingRepository(db *gorm.DB) ILineNotificationSettingRepository {
	return &lineNotificationSettingRepository{db: db}
}

func (r *lineNotificationSettingRepository) FindByUserID(ctx context.Context, userID uint64) (*models.LineNotificationSetting, error) {
	var setting models.LineNotificationSetting
	err := r.db.WithContext(ctx).Preload("User").Where("user_id = ?", userID).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

func (r *lineNotificationSettingRepository) FindAllEnabled(ctx context.Context) ([]models.LineNotificationSetting, error) {
	var settings []models.LineNotificationSetting
	err := r.db.WithContext(ctx).Preload("User").Where("is_enabled = ?", true).Find(&settings).Error
	return settings, err
}

func (r *lineNotificationSettingRepository) Upsert(ctx context.Context, setting *models.LineNotificationSetting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

func (r *lineNotificationSettingRepository) UpdateEnabled(ctx context.Context, userID uint64, isEnabled bool) error {
	return r.db.WithContext(ctx).
		Model(&models.LineNotificationSetting{}).
		Where("user_id = ?", userID).
		Update("is_enabled", isEnabled).Error
}

func (r *lineNotificationSettingRepository) Delete(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.LineNotificationSetting{}).Error
}
