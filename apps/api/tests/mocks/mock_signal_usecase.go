package mocks

import (
	"context"

	"github.com/keeee21/commitly/api/usecase"
)

// MockSignalUsecase is a mock of ISignalUsecase interface.
type MockSignalUsecase struct {
	GetSignalsFunc       func(ctx context.Context, userID uint64, circleID uint64) ([]usecase.Signal, error)
	GetRecentSignalsFunc func(ctx context.Context, userID uint64) ([]usecase.Signal, error)
}

func (m *MockSignalUsecase) GetSignals(ctx context.Context, userID uint64, circleID uint64) ([]usecase.Signal, error) {
	if m.GetSignalsFunc != nil {
		return m.GetSignalsFunc(ctx, userID, circleID)
	}
	return nil, nil
}

func (m *MockSignalUsecase) GetRecentSignals(ctx context.Context, userID uint64) ([]usecase.Signal, error) {
	if m.GetRecentSignalsFunc != nil {
		return m.GetRecentSignalsFunc(ctx, userID)
	}
	return nil, nil
}
