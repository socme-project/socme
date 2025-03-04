package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`
	Logo string
	// Alerts []Alert

	LastAlert time.Time

	ArtemisIsAlive  bool
	ArtemisVersion  string
	ArtemisIP       string
	ArtemisPassword string
	// ...
}

// TODO: We don't want all name, all chars, check IP also
func NewClient(db *gorm.DB, name, logo, artemisIP, artemisPassword string) error {
	client := Client{
		Name:            name,
		Logo:            logo,
		ArtemisIP:       artemisIP,
		ArtemisPassword: artemisPassword,
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
