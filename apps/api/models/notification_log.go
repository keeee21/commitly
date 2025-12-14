package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// NotificationStatus 通知ステータス
type NotificationStatus string

const (
	NotificationStatusSuccess NotificationStatus = "success"
	NotificationStatusFailed  NotificationStatus = "failed"
)

// JSONPayload JSON形式のペイロード
type JSONPayload map[string]interface{}

// Value driver.Valuer インターフェースの実装
func (j JSONPayload) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan sql.Scanner インターフェースの実装
func (j *JSONPayload) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// NotificationLog 通知ログ
type NotificationLog struct {
	ID           uint64             `gorm:"primaryKey;autoIncrement"`
	UserID       uint64             `gorm:"index;not null"`   // FK → users.id
	ChannelType  ChannelType        `gorm:"size:50;not null"` // line / slack / discord
	Period       string             `gorm:"size:20;not null"` // weekly / monthly
	Status       NotificationStatus `gorm:"size:20;not null"` // success / failed
	Payload      JSONPayload        `gorm:"type:jsonb"`       // 送信したメッセージ内容
	ErrorMessage string             `gorm:"type:text"`        // 失敗時のエラーメッセージ
	SentAt       time.Time          `gorm:"not null"`         // 送信日時
	CreatedAt    time.Time          `gorm:"autoCreateTime"`

	// Relations
	User User `gorm:"foreignKey:UserID;references:ID"`
}
