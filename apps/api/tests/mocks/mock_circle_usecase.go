package mocks

import (
	"context"

	"github.com/keeee21/commitly/api/models"
)

// MockCircleUsecase is a mock of ICircleUsecase interface.
type MockCircleUsecase struct {
	GetCirclesFunc   func(ctx context.Context, userID uint64) ([]models.Circle, error)
	CreateCircleFunc func(ctx context.Context, userID uint64, name string) (*models.Circle, error)
	JoinCircleFunc   func(ctx context.Context, userID uint64, inviteCode string) (*models.Circle, error)
	LeaveCircleFunc  func(ctx context.Context, userID uint64, circleID uint64) error
	DeleteCircleFunc func(ctx context.Context, userID uint64, circleID uint64) error
}

func (m *MockCircleUsecase) GetCircles(ctx context.Context, userID uint64) ([]models.Circle, error) {
	if m.GetCirclesFunc != nil {
		return m.GetCirclesFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockCircleUsecase) CreateCircle(ctx context.Context, userID uint64, name string) (*models.Circle, error) {
	if m.CreateCircleFunc != nil {
		return m.CreateCircleFunc(ctx, userID, name)
	}
	return nil, nil
}

func (m *MockCircleUsecase) JoinCircle(ctx context.Context, userID uint64, inviteCode string) (*models.Circle, error) {
	if m.JoinCircleFunc != nil {
		return m.JoinCircleFunc(ctx, userID, inviteCode)
	}
	return nil, nil
}

func (m *MockCircleUsecase) LeaveCircle(ctx context.Context, userID uint64, circleID uint64) error {
	if m.LeaveCircleFunc != nil {
		return m.LeaveCircleFunc(ctx, userID, circleID)
	}
	return nil
}

func (m *MockCircleUsecase) DeleteCircle(ctx context.Context, userID uint64, circleID uint64) error {
	if m.DeleteCircleFunc != nil {
		return m.DeleteCircleFunc(ctx, userID, circleID)
	}
	return nil
}
