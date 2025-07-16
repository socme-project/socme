package model

import (
	"time"

	"github.com/socme-project/backend/utils"
	"gorm.io/gorm"
)

const SESSION_EXPIRATION = 10 * time.Hour

type Session struct {
	Token  string `gorm:"not null;unique;primaryKey"`
	UserID string
	User   User      `gorm:"foreignKey:UserID"`
	Exp    time.Time `gorm:"not null"`
}

func (s Session) String() string {
	return "Session{\n" +
		"\tToken: " + s.Token + "\n" +
		"\tUserID: " + s.UserID + "\n" +
		"\tExp: " + s.Exp.String() + "\n" +
		"}"
}

func GetSession(db *gorm.DB, token string) (*Session, error) {
	s := Session{}

	result := db.Preload("User").First(&s, "token = ?", token)

	if result.Error != nil {
		return nil, result.Error
	}

	return &s, nil
}

func GetSessionByUser(db *gorm.DB, user User) (*Session, error) {
	s := Session{}

	result := db.Preload("User").First(&s, "user_id = ?", user.ID)

	if result.Error != nil {
		return nil, result.Error
	}

	return &s, nil
}

func CreateSession(db *gorm.DB, user User) (*Session, error) {
	token, err := utils.GenerateSecureToken(64)
	if err != nil {
		return nil, err
	}

	existingSession, err := GetSessionByUser(db, user)

	if err != nil {
		s := Session{
			Token:  token,
			UserID: user.ID,
			Exp:    time.Now().Add(SESSION_EXPIRATION),
		}

		result := db.Create(&s)
		if result.Error != nil {
			return nil, result.Error
		}
		return GetSession(db, s.Token)
	} else {
		result := db.Model(&existingSession).Updates(Session{
			Token: token,
			Exp:   time.Now().Add(SESSION_EXPIRATION),
		})

		if result.Error != nil {
			return nil, result.Error
		}
		return GetSession(db, token)
	}
}

func DeleteSession(db *gorm.DB, token string) error {
	_, err := GetSession(db, token)
	if err != nil {
		return err
	}
	result := db.Where("token = ?", token).Delete(&Session{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s Session) IsExpired() bool {
	return time.Now().After(s.Exp)
}
