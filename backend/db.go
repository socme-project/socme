package main

import (
	"backend/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (backend *Backend) initDB() error {
	var err error
	backend.Db, err = gorm.Open(sqlite.Open(backend.DbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	err = backend.Db.AutoMigrate(&model.User{}, &model.Session{}, &model.Client{})
	if err != nil {
		return err
	}

	return nil
}
