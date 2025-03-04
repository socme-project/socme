package model

import "time"

type Session struct {
	ID        uint   `gorm:"primaryKey"`
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
}
