package usecase

import (
	"context"
	"fmt"

	"github.com/keeee21/commitly/api/gateway"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
)

// IRivalUsecase ライバルユースケースのインターフェース
type IRivalUsecase interface {
	GetRivals(ctx context.Context, userID uint64) ([]models.Rival, error)
	AddRival(ctx context.Context, userID uint64, rivalUsername string) (*models.Rival, error)
	RemoveRival(ctx context.Context, userID uint64, rivalID uint64) error
}

type rivalUsecase struct {
	rivalRepo     repository.IRivalRepository
	githubGateway gateway.IGithubGateway
}

// NewRivalUsecase コンストラクタ
func NewRivalUsecase(rivalRepo repository.IRivalRepository, githubGateway gateway.IGithubGateway) IRivalUsecase {
	return &rivalUsecase{
		rivalRepo:     rivalRepo,
		githubGateway: githubGateway,
	}
}

func (u *rivalUsecase) GetRivals(ctx context.Context, userID uint64) ([]models.Rival, error) {
	return u.rivalRepo.FindByUserID(ctx, userID)
}

func (u *rivalUsecase) AddRival(ctx context.Context, userID uint64, rivalUsername string) (*models.Rival, error) {
	// ライバル数の上限チェック
	count, err := u.rivalRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= models.MaxRivalsForFreePlan {
		return nil, fmt.Errorf("ライバル登録数が上限（%d人）に達しています", models.MaxRivalsForFreePlan)
	}

	// Githubユーザー情報を取得
	githubUser, err := u.githubGateway.GetUser(ctx, rivalUsername)
	if err != nil {
		return nil, fmt.Errorf("Githubユーザーが見つかりません: %s", rivalUsername)
	}

	// 既に登録済みかチェック
	exists, err := u.rivalRepo.ExistsByUserIDAndRivalGithubUserID(ctx, userID, githubUser.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("既にライバルとして登録されています: %s", rivalUsername)
	}

	// ライバルを登録
	rival := &models.Rival{
		UserID:              userID,
		RivalGithubUserID:   githubUser.ID,
		RivalGithubUsername: githubUser.Login,
		RivalAvatarURL:      githubUser.AvatarURL,
	}
	if err := u.rivalRepo.Create(ctx, rival); err != nil {
		return nil, err
	}

	return rival, nil
}

func (u *rivalUsecase) RemoveRival(ctx context.Context, userID uint64, rivalID uint64) error {
	// ライバルを取得して所有者を確認
	rival, err := u.rivalRepo.FindByID(ctx, rivalID)
	if err != nil {
		return fmt.Errorf("ライバルが見つかりません")
	}
	if rival.UserID != userID {
		return fmt.Errorf("このライバルを削除する権限がありません")
	}

	return u.rivalRepo.Delete(ctx, rivalID)
}
