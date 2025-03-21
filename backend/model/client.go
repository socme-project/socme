package model

import (
	"fmt"
	"regexp"
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"uniqueIndex"`
	Logo   string
	Alerts []Alert `gorm:"constraint:OnDelete:CASCADE;"` // One-to-many relationship

	LastAlert time.Time

	RefreshRate time.Duration

	WazuhIsAlive    bool
	WazuhVersion    string
	WazuhIP         string
	WazuhPort       string
	WazuhUsername   string
	WazuhPassword   string
	IndexerIP       string
	IndexerPort     string
	IndexerUsername string
	IndexerPassword string
}

func NewClient(
	db *gorm.DB,
	name, logo, wazuhIP, wazuhPort, wazuhUsername, wazuhPassword, indexerIP, indexerPort, indexerUsername, indexerPassword string,
) error {
	if match, _ := regexp.MatchString("^[a-z-]*$", name); !match {
		return fmt.Errorf("Name must be only lowercase letters and -")
	}
	if match, _ := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}$`, wazuhIP); !match {
		return fmt.Errorf("Wazuh IP is not a valid IP")
	}
	if match, _ := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}$`, indexerIP); !match {
		return fmt.Errorf("Wazuh IP is not a valid IP")
	}

	client := Client{
		Name:            name,
		Logo:            logo,
		WazuhIP:         wazuhIP,
		WazuhPort:       wazuhPort,
		WazuhUsername:   wazuhUsername,
		WazuhPassword:   wazuhPassword,
		IndexerIP:       indexerIP,
		IndexerPort:     indexerPort,
		IndexerUsername: indexerUsername,
		IndexerPassword: indexerPassword,
	}

	result := db.Create(&client)
	if result.Error != nil {
		return fmt.Errorf("Error while creating the client.")
	}

	return nil
}

func GetClientByID(db *gorm.DB, id string) (Client, error) {
	var client Client
	if err := db.First(&client, id).Error; err != nil {
		return Client{}, fmt.Errorf("Client not found.")
	}

	return client, nil
}

func GetAllClients(db *gorm.DB) []Client {
	clients := []Client{}
	db.Find(&clients)

	return clients
}
