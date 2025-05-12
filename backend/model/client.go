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
		return fmt.Errorf("Wazuh IP is not a valid IP: %s", wazuhIP)
	}
	if match, _ := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}$`, indexerIP); !match {
		return fmt.Errorf("Indexer IP is not a valid IP: %s", indexerIP)
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
		return fmt.Errorf("Failed to create client: %w", result.Error)
	}

	return nil
}

func GetClientByID(db *gorm.DB, id string) (Client, error) {
	var client Client
	if err := db.First(&client, id).Error; err != nil {
		return Client{}, fmt.Errorf("Client with ID %s not found", id)
	}

	return client, nil
}

func GetAllClients(db *gorm.DB) ([]Client) {
	clients := []Client{}
	_ = db.Find(&clients)
	return clients
}

// Update an existing client
func UpdateClient(db *gorm.DB, id string, name string, logo string,
	wazuhIP string, wazuhPort string, wazuhUsername string, wazuhPassword string,
	indexerIP string, indexerPort string, indexerUsername string, indexerPassword string,
) error {
	// Find the client by ID
	var client Client
	if err := db.First(&client, "id = ?", id).Error; err != nil {
		return err
	}

	// Update client fields
	client.Logo = logo
	client.WazuhIP = wazuhIP
	client.WazuhPort = wazuhPort
	client.WazuhUsername = wazuhUsername
	client.WazuhPassword = wazuhPassword
	client.IndexerIP = indexerIP
	client.IndexerPort = indexerPort
	client.IndexerUsername = indexerUsername
	client.IndexerPassword = indexerPassword

	// Save changes
	return db.Save(&client).Error
}
