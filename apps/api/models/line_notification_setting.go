package models

import "time"

// LineNotificationSetting LINE通知設定
type LineNotificationSetting struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	UserID     uint64    `gorm:"uniqueIndex;not null"` // 1ユーザー1設定
	LineUserID string    `gorm:"size:255;not null"`
	IsEnabled  bool      `gorm:"not null;default:true"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	// Relations
	User User `gorm:"foreignKey:UserID;references:ID"`
}
