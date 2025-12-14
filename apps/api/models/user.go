package models

import "time"

// User ユーザー情報
type User struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement"`
	GithubUserID   uint64    `gorm:"uniqueIndex;not null"`     // Github User ID
	GithubUsername string    `gorm:"size:255;not null"`        // Githubユーザー名
	Email          string    `gorm:"size:255"`                 // メールアドレス
	AvatarURL      string    `gorm:"size:512"`                 // Githubアバター URL
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	// Relations
	Rivals               []Rival               `gorm:"foreignKey:UserID"`
	NotificationSettings []NotificationSetting `gorm:"foreignKey:UserID"`
}
