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

// ActivityItem アクティビティストリームの1件
type ActivityItem struct {
	GithubUsername string `json:"github_username" validate:"required" example:"tanaka"`
	AvatarURL      string `json:"avatar_url" validate:"required" example:"https://avatars.githubusercontent.com/u/1"`
	Repository     string `json:"repository" validate:"required" example:"nextjs-portfolio"`
	CommitCount    int    `json:"commit_count" validate:"required" example:"3"`
	Date           string `json:"date" validate:"required" example:"2026-02-15"`
}

// ActivityStreamResponse アクティビティストリームレスポンス
type ActivityStreamResponse struct {
	Activities []ActivityItem `json:"activities" validate:"required"`
}

// WeeklyRhythm 曜日別コミット有無
type WeeklyRhythm struct {
	Mon bool `json:"mon" example:"true"`
	Tue bool `json:"tue" example:"true"`
	Wed bool `json:"wed" example:"false"`
	Thu bool `json:"thu" example:"true"`
	Fri bool `json:"fri" example:"true"`
	Sat bool `json:"sat" example:"true"`
	Sun bool `json:"sun" example:"true"`
}

// UserRhythm ユーザーのリズム情報
type UserRhythm struct {
	GithubUsername string       `json:"github_username" validate:"required" example:"tanaka"`
	AvatarURL      string       `json:"avatar_url" validate:"required" example:"https://avatars.githubusercontent.com/u/1"`
	PatternLabel   string       `json:"pattern_label" validate:"required" example:"安定型"`
	WeeklyRhythm   WeeklyRhythm `json:"weekly_rhythm" validate:"required"`
}

// RhythmResponse リズム可視化レスポンス
type RhythmResponse struct {
	Users  []UserRhythm `json:"users" validate:"required"`
	Period string       `json:"period" validate:"required" example:"2026-02-09/2026-02-15"`
}
