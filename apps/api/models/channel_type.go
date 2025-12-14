package models

// ChannelType 通知チャンネルタイプ
type ChannelType string

const (
	ChannelTypeLINE    ChannelType = "line"
	ChannelTypeSlack   ChannelType = "slack"
	ChannelTypeDiscord ChannelType = "discord"
)
