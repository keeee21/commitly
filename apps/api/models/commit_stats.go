package models

import "time"

// CommitStats コミット統計（日別・リポジトリ別）
type CommitStats struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement"`
	GithubUserID   uint64    `gorm:"uniqueIndex:idx_commit_stats_unique,priority:1;not null"`           // 対象のGithub User ID
	GithubUsername string    `gorm:"size:255;not null"`                                                 // Githubユーザー名
	Date           time.Time `gorm:"type:date;uniqueIndex:idx_commit_stats_unique,priority:2;not null"` // 日付
	Repository     string    `gorm:"size:255;uniqueIndex:idx_commit_stats_unique,priority:3;not null"`  // リポジトリ名（owner/repo形式）
	CommitCount    int       `gorm:"not null;default:0"`                                                // コミット数
	PrimaryHour    *int      `gorm:"type:smallint"`                                                     // コミットの最頻時間帯（0-23）、nilは未取得
	Language       string    `gorm:"size:100"`                                                          // リポジトリの主要言語
	FetchedAt      time.Time `gorm:"autoCreateTime"`                                                    // 取得日時
}

// TableName テーブル名を指定
func (CommitStats) TableName() string {
	return "commit_stats"
}
