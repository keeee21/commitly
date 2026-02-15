package mocks

import (
	"context"

	"github.com/keeee21/commitly/api/dto"
	"github.com/keeee21/commitly/api/models"
)

// MockActivityUsecase is a mock of IActivityUsecase interface.
type MockActivityUsecase struct {
	GetActivityStreamFunc func(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.ActivityStreamResponse, error)
	GetRhythmFunc         func(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.RhythmResponse, error)
}

func (m *MockActivityUsecase) GetActivityStream(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.ActivityStreamResponse, error) {
	if m.GetActivityStreamFunc != nil {
		return m.GetActivityStreamFunc(ctx, user, rivals)
	}
	return nil, nil
}

func (m *MockActivityUsecase) GetRhythm(ctx context.Context, user *models.User, rivals []models.Rival) (*dto.RhythmResponse, error) {
	if m.GetRhythmFunc != nil {
		return m.GetRhythmFunc(ctx, user, rivals)
	}
	return nil, nil
}
