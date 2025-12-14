package models

import "time"

// DiscordNotificationSetting Discord通知設定
type DiscordNotificationSetting struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	UserID      uint64    `gorm:"uniqueIndex;not null"` // 1ユーザー1設定
	WebhookURL  string    `gorm:"size:512;not null"`
	ServerID    string    `gorm:"size:255"`
	ServerName  string    `gorm:"size:255"`
	ChannelID   string    `gorm:"size:255"`
	ChannelName string    `gorm:"size:255"`
	IsEnabled   bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relations
	User User `gorm:"foreignKey:UserID;references:ID"`
}
