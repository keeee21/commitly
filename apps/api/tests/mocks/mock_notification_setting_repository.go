package mocks

import (
	"context"

	"github.com/keeee21/commitly/api/models"
)

// MockNotificationSettingRepository is a mock of INotificationSettingRepository interface.
type MockNotificationSettingRepository struct {
	FindByUserIDFunc       func(ctx context.Context, userID uint64) ([]models.NotificationSetting, error)
	FindByIDFunc           func(ctx context.Context, id uint64) (*models.NotificationSetting, error)
	FindEnabledSettingsFunc func(ctx context.Context) ([]models.NotificationSetting, error)
	CreateFunc             func(ctx context.Context, setting *models.NotificationSetting) error
	UpdateFunc             func(ctx context.Context, setting *models.NotificationSetting) error
	DeleteFunc             func(ctx context.Context, id uint64) error
}

func (m *MockNotificationSettingRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockNotificationSettingRepository) FindByID(ctx context.Context, id uint64) (*models.NotificationSetting, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockNotificationSettingRepository) FindEnabledSettings(ctx context.Context) ([]models.NotificationSetting, error) {
	if m.FindEnabledSettingsFunc != nil {
		return m.FindEnabledSettingsFunc(ctx)
	}
	return nil, nil
}

func (m *MockNotificationSettingRepository) Create(ctx context.Context, setting *models.NotificationSetting) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, setting)
	}
	return nil
}

func (m *MockNotificationSettingRepository) Update(ctx context.Context, setting *models.NotificationSetting) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, setting)
	}
	return nil
}

func (m *MockNotificationSettingRepository) Delete(ctx context.Context, id uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
