package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	GitHubID string `gorm:"uniqueIndex"`
	Username string
	Avatar   string
	Role     string
}

func GetUserByToken(db *gorm.DB, token string) (User, error) {
	if token == "" {
		return User{}, fmt.Errorf("No token provided")
	}
	var session Session
	if err := db.First(&session, "token = ?", token).Error; err != nil {
		return User{}, fmt.Errorf("Invalid token")
	}

	if time.Now().After(session.ExpiresAt) {
		db.Delete(&session) // Delete expired session
		return User{}, errors.New("session expired")
	}

	var user User
	if err := db.First(&user, session.ID).Error; err != nil {
		return User{}, fmt.Errorf("User not found")
	}

	return user, nil
}

func GetUserByID(db *gorm.DB, userID string) (User, error) {
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return User{}, fmt.Errorf("User not found")
	}

	return user, nil
}

func GetUserFromRequest(db *gorm.DB, c *gin.Context) (User, error) {
	token := c.GetHeader("Authorization")
	return GetUserByToken(db, token)
}

func UpdateUserRole(db *gorm.DB, userID string, role string) error {
	user, err := GetUserByID(db, userID)
	if err != nil {
		return err
	}
	user.Role = role
	db.Save(&user)
	return nil
}
