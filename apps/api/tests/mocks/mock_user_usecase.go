package mocks

import (
	"context"

	"github.com/keeee21/commitly/api/models"
)

// MockUserUsecase is a mock of IUserUsecase interface.
type MockUserUsecase struct {
	GetOrCreateUserFunc       func(ctx context.Context, githubUserID uint64, githubUsername, email, avatarURL string) (*models.User, error)
	GetUserByGithubUserIDFunc func(ctx context.Context, githubUserID uint64) (*models.User, error)
}

func (m *MockUserUsecase) GetOrCreateUser(ctx context.Context, githubUserID uint64, githubUsername, email, avatarURL string) (*models.User, error) {
	if m.GetOrCreateUserFunc != nil {
		return m.GetOrCreateUserFunc(ctx, githubUserID, githubUsername, email, avatarURL)
	}
	return nil, nil
}

func (m *MockUserUsecase) GetUserByGithubUserID(ctx context.Context, githubUserID uint64) (*models.User, error) {
	if m.GetUserByGithubUserIDFunc != nil {
		return m.GetUserByGithubUserIDFunc(ctx, githubUserID)
	}
	return nil, nil
}
