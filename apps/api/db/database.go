package db

import (
	"fmt"

	"github.com/keeee21/commitly/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase creates a new database connection
func NewDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

// AutoMigrate runs auto migration for all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Rival{},
		&models.CommitStats{},
		&models.NotificationSetting{},
	)
}
