package model

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey"`
	GithubID string `gorm:"unique;not null"` // For Github OAuth identification
	Name     string `gorm:"not null"`
	Role     string `gorm:"not null;default:guest"` // e.g., "admin", "user", "guest"
}

// String function for user:
func (u User) String() string {
	return "User{\n" +
		"\tID: " + u.ID + "\n" +
		"\tGithubID: " + u.GithubID + "\n" +
		"\tName: " + u.Name + "\n" +
		"\tRole: " + u.Role + "\n" +
		"}"
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func GetNumberOfUsers(db *gorm.DB) int64 {
	var count int64
	result := db.Model(&User{}).Count(&count)
	if result.Error != nil {
		return 0
	}
	return count
}

func CreateUser(db *gorm.DB, name, githubId string) (*User, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	u := User{
		ID:       uuid.String(),
		Name:     name,
		GithubID: githubId,
	}
	if GetNumberOfUsers(db) == 0 {
		u.Role = "admin" // First user is admin
	}
	result := db.Create(&u)

	if result.Error != nil {
		return nil, result.Error
	}

	return &u, nil
}

func GetUserByID(db *gorm.DB, id string) (*User, error) {
	user := User{}
	result := db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByGithubID(db *gorm.DB, githubID string) (*User, error) {
	user := User{}
	result := db.First(&user, "github_id = ?", githubID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func DeleteUser(db *gorm.DB, id string) error {
	_, err := GetUserByID(db, id)
	if err != nil {
		return err
	}
	result := db.Delete(&User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func EditUserRole(db *gorm.DB, id string, newRole string) error {
	if newRole != "admin" && newRole != "user" && newRole != "guest" {
		return errors.New("invalid role")
	}
	user := User{}
	result := db.First(&user, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	user.Role = newRole
	result = db.Save(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
