package usecase

import (
	"context"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
)

// ISlackNotificationUsecase Slack通知ユースケースのインターフェース
type ISlackNotificationUsecase interface {
	GetSetting(ctx context.Context, userID uint64) (*models.SlackNotificationSetting, error)
	Create(ctx context.Context, userID uint64, webhookURL string) (*models.SlackNotificationSetting, error)
	UpdateEnabled(ctx context.Context, userID uint64, isEnabled bool) error
	Delete(ctx context.Context, userID uint64) error
}

type slackNotificationUsecase struct {
	slackRepo repository.ISlackNotificationSettingRepository
}

// NewSlackNotificationUsecase コンストラクタ
func NewSlackNotificationUsecase(slackRepo repository.ISlackNotificationSettingRepository) ISlackNotificationUsecase {
	return &slackNotificationUsecase{
		slackRepo: slackRepo,
	}
}

func (u *slackNotificationUsecase) GetSetting(ctx context.Context, userID uint64) (*models.SlackNotificationSetting, error) {
	return u.slackRepo.FindByUserID(ctx, userID)
}

func (u *slackNotificationUsecase) Create(ctx context.Context, userID uint64, webhookURL string) (*models.SlackNotificationSetting, error) {
	// 既存の設定を取得
	existing, err := u.slackRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// 更新
		existing.WebhookURL = webhookURL
		existing.IsEnabled = true

		if err := u.slackRepo.Upsert(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// 新規作成
	setting := &models.SlackNotificationSetting{
		UserID:     userID,
		WebhookURL: webhookURL,
		IsEnabled:  true,
	}

	if err := u.slackRepo.Upsert(ctx, setting); err != nil {
		return nil, err
	}

	return setting, nil
}

func (u *slackNotificationUsecase) UpdateEnabled(ctx context.Context, userID uint64, isEnabled bool) error {
	return u.slackRepo.UpdateEnabled(ctx, userID, isEnabled)
}

func (u *slackNotificationUsecase) Delete(ctx context.Context, userID uint64) error {
	return u.slackRepo.Delete(ctx, userID)
}
