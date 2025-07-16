package api

import (
	"github.com/socme-project/backend/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		// TODO: Custom logger
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&model.User{},
		&model.Client{},
		&model.Alert{},
		&model.Session{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
