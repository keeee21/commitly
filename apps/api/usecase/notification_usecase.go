package usecase

import (
	"context"
	"errors"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
)

// INotificationUsecase 通知ユースケースのインターフェース
type INotificationUsecase interface {
	GetSettings(ctx context.Context, userID uint64) ([]models.NotificationSetting, error)
	CreateSetting(ctx context.Context, userID uint64, channelType models.ChannelType, webhookURL, lineUserID string) (*models.NotificationSetting, error)
	UpdateSetting(ctx context.Context, userID uint64, settingID uint64, isEnabled bool, webhookURL, lineUserID string) (*models.NotificationSetting, error)
	DeleteSetting(ctx context.Context, userID uint64, settingID uint64) error
}

type notificationUsecase struct {
	notificationRepo repository.INotificationSettingRepository
}

// NewNotificationUsecase コンストラクタ
func NewNotificationUsecase(notificationRepo repository.INotificationSettingRepository) INotificationUsecase {
	return &notificationUsecase{
		notificationRepo: notificationRepo,
	}
}

func (u *notificationUsecase) GetSettings(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
	return u.notificationRepo.FindByUserID(ctx, userID)
}

func (u *notificationUsecase) CreateSetting(ctx context.Context, userID uint64, channelType models.ChannelType, webhookURL, lineUserID string) (*models.NotificationSetting, error) {
	setting := &models.NotificationSetting{
		UserID:      userID,
		ChannelType: channelType,
		WebhookURL:  webhookURL,
		LINEUserID:  lineUserID,
		IsEnabled:   true,
	}

	if err := u.notificationRepo.Create(ctx, setting); err != nil {
		return nil, err
	}

	return setting, nil
}

func (u *notificationUsecase) UpdateSetting(ctx context.Context, userID uint64, settingID uint64, isEnabled bool, webhookURL, lineUserID string) (*models.NotificationSetting, error) {
	setting, err := u.notificationRepo.FindByID(ctx, settingID)
	if err != nil {
		return nil, err
	}

	if setting.UserID != userID {
		return nil, errors.New("この設定を更新する権限がありません")
	}

	setting.IsEnabled = isEnabled
	setting.WebhookURL = webhookURL
	setting.LINEUserID = lineUserID

	if err := u.notificationRepo.Update(ctx, setting); err != nil {
		return nil, err
	}

	return setting, nil
}

func (u *notificationUsecase) DeleteSetting(ctx context.Context, userID uint64, settingID uint64) error {
	setting, err := u.notificationRepo.FindByID(ctx, settingID)
	if err != nil {
		return err
	}

	if setting.UserID != userID {
		return errors.New("この設定を削除する権限がありません")
	}

	return u.notificationRepo.Delete(ctx, settingID)
}
