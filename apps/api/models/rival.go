package models

import "time"

// Rival ライバル登録
type Rival struct {
	ID                  uint64    `gorm:"primaryKey;autoIncrement"`
	UserID              uint64    `gorm:"index;not null"`    // FK → users.id
	RivalGithubUserID   uint64    `gorm:"not null"`          // ライバルのGithub User ID
	RivalGithubUsername string    `gorm:"size:255;not null"` // ライバルのGithubユーザー名
	RivalAvatarURL      string    `gorm:"size:512"`          // ライバルのGithubアバターURL
	CreatedAt           time.Time `gorm:"autoCreateTime"`

	// Relations
	User User `gorm:"foreignKey:UserID;references:ID"`
}

// MaxRivalsForFreePlan 無料プランのライバル登録上限
const MaxRivalsForFreePlan = 5
