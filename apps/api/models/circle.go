package models

import "time"

// Circle サークル
type Circle struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"size:100;not null"`
	OwnerUserID uint64    `gorm:"index;not null"`
	InviteCode  string    `gorm:"size:32;uniqueIndex;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relations
	Owner   User           `gorm:"foreignKey:OwnerUserID;references:ID"`
	Members []CircleMember `gorm:"foreignKey:CircleID"`
}

// CircleMember サークルメンバー
type CircleMember struct {
	ID       uint64    `gorm:"primaryKey;autoIncrement"`
	CircleID uint64    `gorm:"uniqueIndex:idx_circle_user;not null"`
	UserID   uint64    `gorm:"uniqueIndex:idx_circle_user;not null"`
	JoinedAt time.Time `gorm:"autoCreateTime"`

	// Relations
	Circle Circle `gorm:"foreignKey:CircleID;references:ID"`
	User   User   `gorm:"foreignKey:UserID;references:ID"`
}

// MaxCirclesForFreePlan 無料プランのサークル作成上限
const MaxCirclesForFreePlan = 3

// MaxMembersForFreePlan 無料プランのサークルメンバー上限
const MaxMembersForFreePlan = 3
