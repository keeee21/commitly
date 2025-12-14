package mocks

import (
	"context"

	"github.com/keeee21/commitly/api/models"
)

// MockNotificationUsecase is a mock of INotificationUsecase interface.
type MockNotificationUsecase struct {
	GetSettingsFunc   func(ctx context.Context, userID uint64) ([]models.NotificationSetting, error)
	CreateSettingFunc func(ctx context.Context, userID uint64, channelType models.ChannelType, webhookURL, lineUserID string) (*models.NotificationSetting, error)
	UpdateSettingFunc func(ctx context.Context, userID uint64, settingID uint64, isEnabled bool, webhookURL, lineUserID string) (*models.NotificationSetting, error)
	DeleteSettingFunc func(ctx context.Context, userID uint64, settingID uint64) error
}

func (m *MockNotificationUsecase) GetSettings(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
	if m.GetSettingsFunc != nil {
		return m.GetSettingsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockNotificationUsecase) CreateSetting(ctx context.Context, userID uint64, channelType models.ChannelType, webhookURL, lineUserID string) (*models.NotificationSetting, error) {
	if m.CreateSettingFunc != nil {
		return m.CreateSettingFunc(ctx, userID, channelType, webhookURL, lineUserID)
	}
	return nil, nil
}

func (m *MockNotificationUsecase) UpdateSetting(ctx context.Context, userID uint64, settingID uint64, isEnabled bool, webhookURL, lineUserID string) (*models.NotificationSetting, error) {
	if m.UpdateSettingFunc != nil {
		return m.UpdateSettingFunc(ctx, userID, settingID, isEnabled, webhookURL, lineUserID)
	}
	return nil, nil
}

func (m *MockNotificationUsecase) DeleteSetting(ctx context.Context, userID uint64, settingID uint64) error {
	if m.DeleteSettingFunc != nil {
		return m.DeleteSettingFunc(ctx, userID, settingID)
	}
	return nil
}
