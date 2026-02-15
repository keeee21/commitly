package dto

import "time"

// ErrorResponse エラーレスポンス
type ErrorResponse struct {
	Error string `json:"error" validate:"required" example:"エラーメッセージ"`
}

// MessageResponse メッセージレスポンス
type MessageResponse struct {
	Message string `json:"message" validate:"required" example:"処理が完了しました"`
}

// HealthResponse ヘルスチェックレスポンス
type HealthResponse struct {
	Status string `json:"status" validate:"required" example:"ok"`
}

// UserResponse ユーザーレスポンス
type UserResponse struct {
	ID             uint64    `json:"id" validate:"required" example:"1"`
	GithubUserID   uint64    `json:"github_user_id" validate:"required" example:"12345"`
	GithubUsername string    `json:"github_username" validate:"required" example:"octocat"`
	AvatarURL      string    `json:"avatar_url" validate:"required" example:"https://avatars.githubusercontent.com/u/1"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

// RivalResponse ライバルレスポンス
type RivalResponse struct {
	ID             uint64    `json:"id" validate:"required" example:"1"`
	GithubUserID   uint64    `json:"github_user_id" validate:"required" example:"67890"`
	GithubUsername string    `json:"github_username" validate:"required" example:"rival-user"`
	AvatarURL      string    `json:"avatar_url" validate:"required" example:"https://avatars.githubusercontent.com/u/2"`
	CreatedAt      time.Time `json:"created_at" validate:"required"`
}

// RivalsListResponse ライバル一覧レスポンス
type RivalsListResponse struct {
	Rivals    []RivalResponse `json:"rivals" validate:"required"`
	Count     int             `json:"count" validate:"required" example:"3"`
	MaxRivals int             `json:"max_rivals" validate:"required" example:"5"`
}

// SlackNotificationSettingResponse Slack通知設定レスポンス
type SlackNotificationSettingResponse struct {
	ID         uint64    `json:"id" validate:"required" example:"1"`
	WebhookURL string    `json:"webhook_url" validate:"required" example:"https://hooks.slack.com/services/T00..."`
	IsEnabled  bool      `json:"is_enabled" validate:"required" example:"true"`
	CreatedAt  time.Time `json:"created_at" validate:"required"`
	UpdatedAt  time.Time `json:"updated_at" validate:"required"`
}

// UpdateEnabledResponse 有効/無効更新レスポンス
type UpdateEnabledResponse struct {
	IsEnabled bool `json:"is_enabled" validate:"required" example:"true"`
}
