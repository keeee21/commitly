package mocks

import (
	"context"
	"time"

	"github.com/keeee21/commitly/api/gateway"
)

// MockGithubGateway is a mock of IGithubGateway interface.
type MockGithubGateway struct {
	GetUserFunc              func(ctx context.Context, username string) (*gateway.GithubUser, error)
	GetUserEventsFunc        func(ctx context.Context, username string, page int) ([]gateway.GithubEvent, error)
	GetUserPublicReposFunc   func(ctx context.Context, username string) ([]gateway.GithubRepo, error)
	GetUserContributionsFunc func(ctx context.Context, username string, from, to string) ([]gateway.ContributionDay, error)
	GetRepositoryCommitsFunc func(ctx context.Context, owner, repo, author string, since, until time.Time) ([]gateway.RepositoryCommit, error)
}

func (m *MockGithubGateway) GetUser(ctx context.Context, username string) (*gateway.GithubUser, error) {
	if m.GetUserFunc != nil {
		return m.GetUserFunc(ctx, username)
	}
	return nil, nil
}

func (m *MockGithubGateway) GetUserEvents(ctx context.Context, username string, page int) ([]gateway.GithubEvent, error) {
	if m.GetUserEventsFunc != nil {
		return m.GetUserEventsFunc(ctx, username, page)
	}
	return nil, nil
}

func (m *MockGithubGateway) GetUserPublicRepos(ctx context.Context, username string) ([]gateway.GithubRepo, error) {
	if m.GetUserPublicReposFunc != nil {
		return m.GetUserPublicReposFunc(ctx, username)
	}
	return nil, nil
}

func (m *MockGithubGateway) GetUserContributions(ctx context.Context, username string, from, to string) ([]gateway.ContributionDay, error) {
	if m.GetUserContributionsFunc != nil {
		return m.GetUserContributionsFunc(ctx, username, from, to)
	}
	return nil, nil
}

func (m *MockGithubGateway) GetRepositoryCommits(ctx context.Context, owner, repo, author string, since, until time.Time) ([]gateway.RepositoryCommit, error) {
	if m.GetRepositoryCommitsFunc != nil {
		return m.GetRepositoryCommitsFunc(ctx, owner, repo, author, since, until)
	}
	return nil, nil
}
