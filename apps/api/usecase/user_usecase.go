package usecase

import (
	"context"

	"github.com/keeee21/commitly/api/gateway"
	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
)

// IUserUsecase ユーザーユースケースのインターフェース
type IUserUsecase interface {
	GetOrCreateUser(ctx context.Context, githubUserID uint64, githubUsername, email, avatarURL string) (*models.User, error)
	GetUserByGithubUserID(ctx context.Context, githubUserID uint64) (*models.User, error)
}

type userUsecase struct {
	userRepo      repository.IUserRepository
	githubGateway gateway.IGithubGateway
}

// NewUserUsecase コンストラクタ
func NewUserUsecase(userRepo repository.IUserRepository, githubGateway gateway.IGithubGateway) IUserUsecase {
	return &userUsecase{
		userRepo:      userRepo,
		githubGateway: githubGateway,
	}
}

func (u *userUsecase) GetOrCreateUser(ctx context.Context, githubUserID uint64, githubUsername, email, avatarURL string) (*models.User, error) {
	// 既存ユーザーを検索
	user, err := u.userRepo.FindByGithubUserID(ctx, githubUserID)
	if err == nil {
		// ユーザーが存在する場合、情報を更新
		user.GithubUsername = githubUsername
		user.Email = email
		user.AvatarURL = avatarURL
		if err := u.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// 新規ユーザーを作成
	newUser := &models.User{
		GithubUserID:   githubUserID,
		GithubUsername: githubUsername,
		Email:          email,
		AvatarURL:      avatarURL,
	}
	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (u *userUsecase) GetUserByGithubUserID(ctx context.Context, githubUserID uint64) (*models.User, error) {
	return u.userRepo.FindByGithubUserID(ctx, githubUserID)
}
