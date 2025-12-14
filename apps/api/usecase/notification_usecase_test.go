package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/keeee21/commitly/api/models"
	"github.com/stretchr/testify/assert"
)

// mockNotificationSettingRepository テスト用のモックリポジトリ
type mockNotificationSettingRepository struct {
	FindByUserIDFunc func(ctx context.Context, userID uint64) ([]models.NotificationSetting, error)
	FindByIDFunc     func(ctx context.Context, id uint64) (*models.NotificationSetting, error)
	CreateFunc       func(ctx context.Context, setting *models.NotificationSetting) error
	UpdateFunc       func(ctx context.Context, setting *models.NotificationSetting) error
	DeleteFunc       func(ctx context.Context, id uint64) error
}

func (m *mockNotificationSettingRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockNotificationSettingRepository) FindByID(ctx context.Context, id uint64) (*models.NotificationSetting, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockNotificationSettingRepository) Create(ctx context.Context, setting *models.NotificationSetting) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, setting)
	}
	return nil
}

func (m *mockNotificationSettingRepository) Update(ctx context.Context, setting *models.NotificationSetting) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, setting)
	}
	return nil
}

func (m *mockNotificationSettingRepository) Delete(ctx context.Context, id uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *mockNotificationSettingRepository) FindEnabledSettings(ctx context.Context) ([]models.NotificationSetting, error) {
	return nil, nil
}

func TestGetSettings_Success(t *testing.T) {
	ctx := context.Background()

	expectedSettings := []models.NotificationSetting{
		{ID: 1, UserID: 1, ChannelType: models.ChannelTypeLINE, IsEnabled: true},
		{ID: 2, UserID: 1, ChannelType: models.ChannelTypeSlack, WebhookURL: "https://hooks.slack.com/xxx", IsEnabled: true},
	}

	mockRepo := &mockNotificationSettingRepository{
		FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
			assert.Equal(t, uint64(1), userID)
			return expectedSettings, nil
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	settings, err := usecase.GetSettings(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, settings, 2)
	assert.Equal(t, models.ChannelTypeLINE, settings[0].ChannelType)
	assert.Equal(t, models.ChannelTypeSlack, settings[1].ChannelType)
}

func TestGetSettings_Empty(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockNotificationSettingRepository{
		FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
			return []models.NotificationSetting{}, nil
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	settings, err := usecase.GetSettings(ctx, 1)

	assert.NoError(t, err)
	assert.Empty(t, settings)
}

func TestGetSettings_Error(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockNotificationSettingRepository{
		FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.NotificationSetting, error) {
			return nil, errors.New("database error")
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	settings, err := usecase.GetSettings(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, settings)
}

func TestCreateSetting_Success(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockNotificationSettingRepository{
		CreateFunc: func(ctx context.Context, setting *models.NotificationSetting) error {
			setting.ID = 1
			return nil
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	setting, err := usecase.CreateSetting(ctx, 1, models.ChannelTypeSlack, "https://hooks.slack.com/xxx", "")

	assert.NoError(t, err)
	assert.NotNil(t, setting)
	assert.Equal(t, uint64(1), setting.ID)
	assert.Equal(t, uint64(1), setting.UserID)
	assert.Equal(t, models.ChannelTypeSlack, setting.ChannelType)
	assert.Equal(t, "https://hooks.slack.com/xxx", setting.WebhookURL)
	assert.True(t, setting.IsEnabled)
}

func TestCreateSetting_LINE(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockNotificationSettingRepository{
		CreateFunc: func(ctx context.Context, setting *models.NotificationSetting) error {
			setting.ID = 1
			return nil
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	setting, err := usecase.CreateSetting(ctx, 1, models.ChannelTypeLINE, "", "U1234567890")

	assert.NoError(t, err)
	assert.NotNil(t, setting)
	assert.Equal(t, models.ChannelTypeLINE, setting.ChannelType)
	assert.Equal(t, "U1234567890", setting.LINEUserID)
}

func TestCreateSetting_Error(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockNotificationSettingRepository{
		CreateFunc: func(ctx context.Context, setting *models.NotificationSetting) error {
			return errors.New("database error")
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	setting, err := usecase.CreateSetting(ctx, 1, models.ChannelTypeSlack, "https://hooks.slack.com/xxx", "")

	assert.Error(t, err)
	assert.Nil(t, setting)
}

func TestUpdateSetting_Success(t *testing.T) {
	ctx := context.Background()

	existingSetting := &models.NotificationSetting{
		ID:          1,
		UserID:      1,
		ChannelType: models.ChannelTypeSlack,
		WebhookURL:  "https://old.url",
		IsEnabled:   true,
	}

	mockRepo := &mockNotificationSettingRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.NotificationSetting, error) {
			return existingSetting, nil
		},
		UpdateFunc: func(ctx context.Context, setting *models.NotificationSetting) error {
			return nil
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	setting, err := usecase.UpdateSetting(ctx, 1, 1, false, "https://new.url", "")

	assert.NoError(t, err)
	assert.NotNil(t, setting)
	assert.False(t, setting.IsEnabled)
	assert.Equal(t, "https://new.url", setting.WebhookURL)
}

func TestUpdateSetting_WrongUser(t *testing.T) {
	ctx := context.Background()

	existingSetting := &models.NotificationSetting{
		ID:     1,
		UserID: 999, // Different user
	}

	mockRepo := &mockNotificationSettingRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.NotificationSetting, error) {
			return existingSetting, nil
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	setting, err := usecase.UpdateSetting(ctx, 1, 1, false, "", "")

	assert.Error(t, err)
	assert.Nil(t, setting)
}

func TestUpdateSetting_NotFound(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockNotificationSettingRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.NotificationSetting, error) {
			return nil, errors.New("not found")
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	setting, err := usecase.UpdateSetting(ctx, 1, 1, false, "", "")

	assert.Error(t, err)
	assert.Nil(t, setting)
}

func TestDeleteSetting_Success(t *testing.T) {
	ctx := context.Background()

	existingSetting := &models.NotificationSetting{
		ID:     1,
		UserID: 1,
	}

	mockRepo := &mockNotificationSettingRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.NotificationSetting, error) {
			return existingSetting, nil
		},
		DeleteFunc: func(ctx context.Context, id uint64) error {
			return nil
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	err := usecase.DeleteSetting(ctx, 1, 1)

	assert.NoError(t, err)
}

func TestDeleteSetting_WrongUser(t *testing.T) {
	ctx := context.Background()

	existingSetting := &models.NotificationSetting{
		ID:     1,
		UserID: 999, // Different user
	}

	mockRepo := &mockNotificationSettingRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.NotificationSetting, error) {
			return existingSetting, nil
		},
	}

	usecase := NewNotificationUsecase(mockRepo)
	err := usecase.DeleteSetting(ctx, 1, 1)

	assert.Error(t, err)
}
