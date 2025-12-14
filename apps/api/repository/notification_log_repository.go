package repository

import (
	"context"
	"time"

	"github.com/keeee21/commitly/api/models"
	"gorm.io/gorm"
)

// INotificationLogRepository 通知ログリポジトリのインターフェース
type INotificationLogRepository interface {
	Create(ctx context.Context, log *models.NotificationLog) error
	FindByUserID(ctx context.Context, userID uint64, limit int) ([]models.NotificationLog, error)
	FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NotificationLog, error)
	FindByPeriodAndDate(ctx context.Context, period string, date time.Time) ([]models.NotificationLog, error)
}

type notificationLogRepository struct {
	db *gorm.DB
}

// NewNotificationLogRepository コンストラクタ
func NewNotificationLogRepository(db *gorm.DB) INotificationLogRepository {
	return &notificationLogRepository{db: db}
}

func (r *notificationLogRepository) Create(ctx context.Context, log *models.NotificationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *notificationLogRepository) FindByUserID(ctx context.Context, userID uint64, limit int) ([]models.NotificationLog, error) {
	var logs []models.NotificationLog
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("sent_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *notificationLogRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NotificationLog, error) {
	var logs []models.NotificationLog
	if err := r.db.WithContext(ctx).
		Where("sent_at >= ? AND sent_at <= ?", startDate, endDate).
		Order("sent_at DESC").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *notificationLogRepository) FindByPeriodAndDate(ctx context.Context, period string, date time.Time) ([]models.NotificationLog, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var logs []models.NotificationLog
	if err := r.db.WithContext(ctx).
		Where("period = ? AND sent_at >= ? AND sent_at < ?", period, startOfDay, endOfDay).
		Order("sent_at DESC").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
