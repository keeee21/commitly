package usecase

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/keeee21/commitly/api/models"
	"github.com/keeee21/commitly/api/repository"
)

// ICircleUsecase サークルユースケースのインターフェース
type ICircleUsecase interface {
	GetCircles(ctx context.Context, userID uint64) ([]models.Circle, error)
	CreateCircle(ctx context.Context, userID uint64, name string) (*models.Circle, error)
	JoinCircle(ctx context.Context, userID uint64, inviteCode string) (*models.Circle, error)
	LeaveCircle(ctx context.Context, userID uint64, circleID uint64) error
	DeleteCircle(ctx context.Context, userID uint64, circleID uint64) error
}

type circleUsecase struct {
	circleRepo repository.ICircleRepository
}

// NewCircleUsecase コンストラクタ
func NewCircleUsecase(circleRepo repository.ICircleRepository) ICircleUsecase {
	return &circleUsecase{
		circleRepo: circleRepo,
	}
}

func (u *circleUsecase) GetCircles(ctx context.Context, userID uint64) ([]models.Circle, error) {
	return u.circleRepo.FindByUserID(ctx, userID)
}

func (u *circleUsecase) CreateCircle(ctx context.Context, userID uint64, name string) (*models.Circle, error) {
	count, err := u.circleRepo.CountByOwnerUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if count >= models.MaxCirclesForFreePlan {
		return nil, fmt.Errorf("サークル作成数が上限（%d個）に達しています", models.MaxCirclesForFreePlan)
	}

	inviteCode, err := generateInviteCode()
	if err != nil {
		return nil, fmt.Errorf("招待コードの生成に失敗しました")
	}

	circle := &models.Circle{
		Name:        name,
		OwnerUserID: userID,
		InviteCode:  inviteCode,
	}
	if err := u.circleRepo.Create(ctx, circle); err != nil {
		return nil, err
	}

	member := &models.CircleMember{
		CircleID: circle.ID,
		UserID:   userID,
	}
	if err := u.circleRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}

	return u.circleRepo.FindByID(ctx, circle.ID)
}

func (u *circleUsecase) JoinCircle(ctx context.Context, userID uint64, inviteCode string) (*models.Circle, error) {
	circle, err := u.circleRepo.FindByInviteCode(ctx, inviteCode)
	if err != nil {
		return nil, fmt.Errorf("招待コードが無効です")
	}

	isMember, err := u.circleRepo.IsMember(ctx, circle.ID, userID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, fmt.Errorf("既にこのサークルに参加しています")
	}

	memberCount, err := u.circleRepo.CountMembers(ctx, circle.ID)
	if err != nil {
		return nil, err
	}
	if memberCount >= models.MaxMembersForFreePlan {
		return nil, fmt.Errorf("サークルのメンバー数が上限（%d人）に達しています", models.MaxMembersForFreePlan)
	}

	member := &models.CircleMember{
		CircleID: circle.ID,
		UserID:   userID,
	}
	if err := u.circleRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}

	return u.circleRepo.FindByID(ctx, circle.ID)
}

func (u *circleUsecase) LeaveCircle(ctx context.Context, userID uint64, circleID uint64) error {
	circle, err := u.circleRepo.FindByID(ctx, circleID)
	if err != nil {
		return fmt.Errorf("サークルが見つかりません")
	}

	if circle.OwnerUserID == userID {
		return fmt.Errorf("オーナーはサークルを退会できません")
	}

	isMember, err := u.circleRepo.IsMember(ctx, circleID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return fmt.Errorf("このサークルのメンバーではありません")
	}

	return u.circleRepo.RemoveMember(ctx, circleID, userID)
}

func (u *circleUsecase) DeleteCircle(ctx context.Context, userID uint64, circleID uint64) error {
	circle, err := u.circleRepo.FindByID(ctx, circleID)
	if err != nil {
		return fmt.Errorf("サークルが見つかりません")
	}

	if circle.OwnerUserID != userID {
		return fmt.Errorf("サークルを削除する権限がありません")
	}

	return u.circleRepo.Delete(ctx, circleID)
}

func generateInviteCode() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}
