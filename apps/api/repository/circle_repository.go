package repository

import (
	"context"

	"github.com/keeee21/commitly/api/models"
	"gorm.io/gorm"
)

// ICircleRepository サークルリポジトリのインターフェース
type ICircleRepository interface {
	FindByUserID(ctx context.Context, userID uint64) ([]models.Circle, error)
	FindByID(ctx context.Context, id uint64) (*models.Circle, error)
	FindByInviteCode(ctx context.Context, code string) (*models.Circle, error)
	CountByOwnerUserID(ctx context.Context, userID uint64) (int64, error)
	Create(ctx context.Context, circle *models.Circle) error
	Delete(ctx context.Context, id uint64) error
	AddMember(ctx context.Context, member *models.CircleMember) error
	RemoveMember(ctx context.Context, circleID, userID uint64) error
	CountMembers(ctx context.Context, circleID uint64) (int64, error)
	IsMember(ctx context.Context, circleID, userID uint64) (bool, error)
}

type circleRepository struct {
	db *gorm.DB
}

// NewCircleRepository コンストラクタ
func NewCircleRepository(db *gorm.DB) ICircleRepository {
	return &circleRepository{db: db}
}

func (r *circleRepository) FindByUserID(ctx context.Context, userID uint64) ([]models.Circle, error) {
	var circles []models.Circle
	if err := r.db.WithContext(ctx).
		Joins("JOIN circle_members ON circle_members.circle_id = circles.id").
		Where("circle_members.user_id = ?", userID).
		Preload("Members").Preload("Members.User").
		Find(&circles).Error; err != nil {
		return nil, err
	}
	return circles, nil
}

func (r *circleRepository) FindByID(ctx context.Context, id uint64) (*models.Circle, error) {
	var circle models.Circle
	if err := r.db.WithContext(ctx).
		Preload("Members").Preload("Members.User").
		First(&circle, id).Error; err != nil {
		return nil, err
	}
	return &circle, nil
}

func (r *circleRepository) FindByInviteCode(ctx context.Context, code string) (*models.Circle, error) {
	var circle models.Circle
	if err := r.db.WithContext(ctx).
		Preload("Members").Preload("Members.User").
		Where("invite_code = ?", code).
		First(&circle).Error; err != nil {
		return nil, err
	}
	return &circle, nil
}

func (r *circleRepository) CountByOwnerUserID(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Circle{}).Where("owner_user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *circleRepository) Create(ctx context.Context, circle *models.Circle) error {
	return r.db.WithContext(ctx).Create(circle).Error
}

func (r *circleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Select("Members").Delete(&models.Circle{ID: id}).Error
}

func (r *circleRepository) AddMember(ctx context.Context, member *models.CircleMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *circleRepository) RemoveMember(ctx context.Context, circleID, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("circle_id = ? AND user_id = ?", circleID, userID).
		Delete(&models.CircleMember{}).Error
}

func (r *circleRepository) CountMembers(ctx context.Context, circleID uint64) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.CircleMember{}).Where("circle_id = ?", circleID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *circleRepository) IsMember(ctx context.Context, circleID, userID uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.CircleMember{}).
		Where("circle_id = ? AND user_id = ?", circleID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
