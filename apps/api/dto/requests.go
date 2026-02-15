package dto

// CallbackRequest Github OAuth コールバックリクエスト
type CallbackRequest struct {
	GithubUserID   uint64 `json:"github_user_id" validate:"required"`
	GithubUsername string `json:"github_username" validate:"required"`
	Email          string `json:"email"`
	AvatarURL      string `json:"avatar_url"`
}

// AddRivalRequest ライバル追加リクエスト
type AddRivalRequest struct {
	Username string `json:"username" validate:"required"`
}

// CreateSlackNotificationRequest Slack通知設定作成リクエスト
type CreateSlackNotificationRequest struct {
	WebhookURL string `json:"webhook_url"`
}

// UpdateEnabledRequest 有効/無効更新リクエスト
type UpdateEnabledRequest struct {
	IsEnabled bool `json:"is_enabled"`
}

// CreateCircleRequest サークル作成リクエスト
type CreateCircleRequest struct {
	Name string `json:"name" validate:"required"`
}

// JoinCircleRequest サークル参加リクエスト
type JoinCircleRequest struct {
	InviteCode string `json:"invite_code" validate:"required"`
}
