package models

import "time"

// ChannelType 通知チャンネルタイプ
type ChannelType string

const (
	ChannelTypeLINE    ChannelType = "line"
	ChannelTypeSlack   ChannelType = "slack"
	ChannelTypeDiscord ChannelType = "discord"
)

// NotificationSetting 通知設定
type NotificationSetting struct {
	ID          uint64      `gorm:"primaryKey;autoIncrement"`
	UserID      uint64      `gorm:"index;not null"`           // FK → users.id
	ChannelType ChannelType `gorm:"size:50;not null"`         // line / slack / discord
	WebhookURL  string      `gorm:"size:512"`                 // Webhook URL（Slack/Discord）
	LINEUserID  string      `gorm:"size:255"`                 // LINE User ID
	IsEnabled   bool        `gorm:"not null;default:true"`    // 有効/無効
	CreatedAt   time.Time   `gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime"`

	// Relations
	User User `gorm:"foreignKey:UserID;references:ID"`
}
