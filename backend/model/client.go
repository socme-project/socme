package model

import (
	"fmt"
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

// TODO: We don't want all name, all chars, check IP also
func NewClient(
	db *gorm.DB,
	name, logo, wazuhIP, wazuhPort, wazuhUsername, wazuhPassword, indexerIP, indexerPort, indexerUsername, indexerPassword string,
) error {
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
		return fmt.Errorf("Error while creating the client")
	}

	return nil
}

func GetClientByName(db *gorm.DB, name string) (Client, error) {
	var client Client
	if err := db.First(&client, "name = ?", name).Error; err != nil {
		return Client{}, fmt.Errorf("Client not found")
	}

	return client, nil
}

func GetClientByID(db *gorm.DB, id string) (Client, error) {
	var client Client
	if err := db.First(&client, id).Error; err != nil {
		return Client{}, fmt.Errorf("Client not found")
	}

	return client, nil
}

func GetAllClients(db *gorm.DB) []Client {
	clients := []Client{}
	db.Find(&clients)

	return clients
}
