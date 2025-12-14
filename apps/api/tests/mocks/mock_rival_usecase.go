package mocks

import (
	"context"

	"github.com/keeee21/commitly/api/models"
)

// MockRivalUsecase is a mock of IRivalUsecase interface.
type MockRivalUsecase struct {
	GetRivalsFunc   func(ctx context.Context, userID uint64) ([]models.Rival, error)
	AddRivalFunc    func(ctx context.Context, userID uint64, rivalUsername string) (*models.Rival, error)
	RemoveRivalFunc func(ctx context.Context, userID uint64, rivalID uint64) error
}

func (m *MockRivalUsecase) GetRivals(ctx context.Context, userID uint64) ([]models.Rival, error) {
	if m.GetRivalsFunc != nil {
		return m.GetRivalsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockRivalUsecase) AddRival(ctx context.Context, userID uint64, rivalUsername string) (*models.Rival, error) {
	if m.AddRivalFunc != nil {
		return m.AddRivalFunc(ctx, userID, rivalUsername)
	}
	return nil, nil
}

func (m *MockRivalUsecase) RemoveRival(ctx context.Context, userID uint64, rivalID uint64) error {
	if m.RemoveRivalFunc != nil {
		return m.RemoveRivalFunc(ctx, userID, rivalID)
	}
	return nil
}
