package repository

import (
	"context"

	"github.com/keeee21/commitly/api/models"
	"gorm.io/gorm"
)

// IUserRepository ユーザーリポジトリのインターフェース
type IUserRepository interface {
	FindByID(ctx context.Context, id uint64) (*models.User, error)
	FindByGithubUserID(ctx context.Context, githubUserID uint64) (*models.User, error)
	FindAll(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository コンストラクタ
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByGithubUserID(ctx context.Context, githubUserID uint64) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("github_user_id = ?", githubUserID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
