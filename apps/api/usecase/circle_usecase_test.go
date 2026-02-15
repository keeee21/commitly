package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/keeee21/commitly/api/models"
	"github.com/stretchr/testify/assert"
)

type circleMockCircleRepository struct {
	FindByUserIDFunc       func(ctx context.Context, userID uint64) ([]models.Circle, error)
	FindByIDFunc           func(ctx context.Context, id uint64) (*models.Circle, error)
	FindByInviteCodeFunc   func(ctx context.Context, code string) (*models.Circle, error)
	CountByOwnerUserIDFunc func(ctx context.Context, userID uint64) (int64, error)
	CreateFunc             func(ctx context.Context, circle *models.Circle) error
	DeleteFunc             func(ctx context.Context, id uint64) error
	AddMemberFunc          func(ctx context.Context, member *models.CircleMember) error
	RemoveMemberFunc       func(ctx context.Context, circleID, userID uint64) error
	CountMembersFunc       func(ctx context.Context, circleID uint64) (int64, error)
	IsMemberFunc           func(ctx context.Context, circleID, userID uint64) (bool, error)
}

func (m *circleMockCircleRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.Circle, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *circleMockCircleRepository) FindByID(ctx context.Context, id uint64) (*models.Circle, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *circleMockCircleRepository) FindByInviteCode(ctx context.Context, code string) (*models.Circle, error) {
	if m.FindByInviteCodeFunc != nil {
		return m.FindByInviteCodeFunc(ctx, code)
	}
	return nil, nil
}

func (m *circleMockCircleRepository) CountByOwnerUserID(ctx context.Context, userID uint64) (int64, error) {
	if m.CountByOwnerUserIDFunc != nil {
		return m.CountByOwnerUserIDFunc(ctx, userID)
	}
	return 0, nil
}

func (m *circleMockCircleRepository) Create(ctx context.Context, circle *models.Circle) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, circle)
	}
	return nil
}

func (m *circleMockCircleRepository) Delete(ctx context.Context, id uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *circleMockCircleRepository) AddMember(ctx context.Context, member *models.CircleMember) error {
	if m.AddMemberFunc != nil {
		return m.AddMemberFunc(ctx, member)
	}
	return nil
}

func (m *circleMockCircleRepository) RemoveMember(ctx context.Context, circleID, userID uint64) error {
	if m.RemoveMemberFunc != nil {
		return m.RemoveMemberFunc(ctx, circleID, userID)
	}
	return nil
}

func (m *circleMockCircleRepository) CountMembers(ctx context.Context, circleID uint64) (int64, error) {
	if m.CountMembersFunc != nil {
		return m.CountMembersFunc(ctx, circleID)
	}
	return 0, nil
}

func (m *circleMockCircleRepository) IsMember(ctx context.Context, circleID, userID uint64) (bool, error) {
	if m.IsMemberFunc != nil {
		return m.IsMemberFunc(ctx, circleID, userID)
	}
	return false, nil
}

func TestGetCircles_Success(t *testing.T) {
	ctx := context.Background()
	expected := []models.Circle{
		{ID: 1, Name: "Circle1", OwnerUserID: 1},
		{ID: 2, Name: "Circle2", OwnerUserID: 2},
	}

	mockRepo := &circleMockCircleRepository{
		FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.Circle, error) {
			return expected, nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	circles, err := uc.GetCircles(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, circles, 2)
	assert.Equal(t, expected, circles)
}

func TestGetCircles_Empty(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		FindByUserIDFunc: func(ctx context.Context, userID uint64) ([]models.Circle, error) {
			return []models.Circle{}, nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	circles, err := uc.GetCircles(ctx, 1)

	assert.NoError(t, err)
	assert.Empty(t, circles)
}

func TestCreateCircle_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		CountByOwnerUserIDFunc: func(ctx context.Context, userID uint64) (int64, error) {
			return 0, nil
		},
		CreateFunc: func(ctx context.Context, circle *models.Circle) error {
			circle.ID = 1
			return nil
		},
		AddMemberFunc: func(ctx context.Context, member *models.CircleMember) error {
			return nil
		},
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return &models.Circle{
				ID:          1,
				Name:        "テストサークル",
				OwnerUserID: 1,
				InviteCode:  "abcd1234",
				CreatedAt:   time.Now(),
				Members: []models.CircleMember{
					{ID: 1, CircleID: 1, UserID: 1},
				},
			}, nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	circle, err := uc.CreateCircle(ctx, 1, "テストサークル")

	assert.NoError(t, err)
	assert.NotNil(t, circle)
	assert.Equal(t, "テストサークル", circle.Name)
	assert.Equal(t, uint64(1), circle.OwnerUserID)
}

func TestCreateCircle_MaxReached(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		CountByOwnerUserIDFunc: func(ctx context.Context, userID uint64) (int64, error) {
			return int64(models.MaxCirclesForFreePlan), nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	circle, err := uc.CreateCircle(ctx, 1, "テストサークル")

	assert.Error(t, err)
	assert.Nil(t, circle)
	assert.Contains(t, err.Error(), "上限")
}

func TestJoinCircle_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		FindByInviteCodeFunc: func(ctx context.Context, code string) (*models.Circle, error) {
			return &models.Circle{ID: 1, Name: "Circle1", OwnerUserID: 2, InviteCode: "abcd1234"}, nil
		},
		IsMemberFunc: func(ctx context.Context, circleID, userID uint64) (bool, error) {
			return false, nil
		},
		CountMembersFunc: func(ctx context.Context, circleID uint64) (int64, error) {
			return 1, nil
		},
		AddMemberFunc: func(ctx context.Context, member *models.CircleMember) error {
			return nil
		},
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return &models.Circle{
				ID: 1, Name: "Circle1", OwnerUserID: 2, InviteCode: "abcd1234",
				Members: []models.CircleMember{
					{ID: 1, CircleID: 1, UserID: 2},
					{ID: 2, CircleID: 1, UserID: 1},
				},
			}, nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	circle, err := uc.JoinCircle(ctx, 1, "abcd1234")

	assert.NoError(t, err)
	assert.NotNil(t, circle)
	assert.Len(t, circle.Members, 2)
}

func TestJoinCircle_InvalidCode(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		FindByInviteCodeFunc: func(ctx context.Context, code string) (*models.Circle, error) {
			return nil, errors.New("not found")
		},
	}

	uc := NewCircleUsecase(mockRepo)
	circle, err := uc.JoinCircle(ctx, 1, "invalid")

	assert.Error(t, err)
	assert.Nil(t, circle)
	assert.Contains(t, err.Error(), "招待コードが無効です")
}

func TestJoinCircle_AlreadyMember(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		FindByInviteCodeFunc: func(ctx context.Context, code string) (*models.Circle, error) {
			return &models.Circle{ID: 1, Name: "Circle1", OwnerUserID: 2}, nil
		},
		IsMemberFunc: func(ctx context.Context, circleID, userID uint64) (bool, error) {
			return true, nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	circle, err := uc.JoinCircle(ctx, 1, "abcd1234")

	assert.Error(t, err)
	assert.Nil(t, circle)
	assert.Contains(t, err.Error(), "既にこのサークルに参加しています")
}

func TestJoinCircle_MaxMembers(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		FindByInviteCodeFunc: func(ctx context.Context, code string) (*models.Circle, error) {
			return &models.Circle{ID: 1, Name: "Circle1", OwnerUserID: 2}, nil
		},
		IsMemberFunc: func(ctx context.Context, circleID, userID uint64) (bool, error) {
			return false, nil
		},
		CountMembersFunc: func(ctx context.Context, circleID uint64) (int64, error) {
			return int64(models.MaxMembersForFreePlan), nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	circle, err := uc.JoinCircle(ctx, 1, "abcd1234")

	assert.Error(t, err)
	assert.Nil(t, circle)
	assert.Contains(t, err.Error(), "上限")
}

func TestLeaveCircle_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return &models.Circle{ID: 1, OwnerUserID: 2}, nil
		},
		IsMemberFunc: func(ctx context.Context, circleID, userID uint64) (bool, error) {
			return true, nil
		},
		RemoveMemberFunc: func(ctx context.Context, circleID, userID uint64) error {
			return nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	err := uc.LeaveCircle(ctx, 1, 1)

	assert.NoError(t, err)
}

func TestLeaveCircle_OwnerCannotLeave(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return &models.Circle{ID: 1, OwnerUserID: 1}, nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	err := uc.LeaveCircle(ctx, 1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "オーナーはサークルを退会できません")
}

func TestDeleteCircle_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return &models.Circle{ID: 1, OwnerUserID: 1}, nil
		},
		DeleteFunc: func(ctx context.Context, id uint64) error {
			return nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	err := uc.DeleteCircle(ctx, 1, 1)

	assert.NoError(t, err)
}

func TestDeleteCircle_NotOwner(t *testing.T) {
	ctx := context.Background()
	mockRepo := &circleMockCircleRepository{
		FindByIDFunc: func(ctx context.Context, id uint64) (*models.Circle, error) {
			return &models.Circle{ID: 1, OwnerUserID: 2}, nil
		},
	}

	uc := NewCircleUsecase(mockRepo)
	err := uc.DeleteCircle(ctx, 1, 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "権限がありません")
}
