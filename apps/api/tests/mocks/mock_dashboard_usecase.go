package mocks

import (
	"context"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/usecase"
)

// MockDashboardUsecase is a mock of IDashboardUsecase interface.
type MockDashboardUsecase struct {
	GetWeeklyDashboardFunc  func(ctx context.Context, user *models.User, rivals []models.Rival) (*usecase.DashboardData, error)
	GetMonthlyDashboardFunc func(ctx context.Context, user *models.User, rivals []models.Rival) (*usecase.DashboardData, error)
}

func (m *MockDashboardUsecase) GetWeeklyDashboard(ctx context.Context, user *models.User, rivals []models.Rival) (*usecase.DashboardData, error) {
	if m.GetWeeklyDashboardFunc != nil {
		return m.GetWeeklyDashboardFunc(ctx, user, rivals)
	}
	return nil, nil
}

func (m *MockDashboardUsecase) GetMonthlyDashboard(ctx context.Context, user *models.User, rivals []models.Rival) (*usecase.DashboardData, error) {
	if m.GetMonthlyDashboardFunc != nil {
		return m.GetMonthlyDashboardFunc(ctx, user, rivals)
	}
	return nil, nil
}
