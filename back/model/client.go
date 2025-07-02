package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/socme-projects/backend/utils"
	"gorm.io/gorm"
)

type Client struct {
	ID     string `gorm:"primaryKey"`
	Name   string `gorm:"unique;not null"`
	Logo   string
	Alerts []Alert `gorm:"constraint:OnDelete:CASCADE;"`

	LastAlert time.Time

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

func (c Client) String() string {
	return "Client{\n" +
		"\tID: " + c.ID + "\n" +
		"\tName: " + c.Name + "\n" +
		"\tLogo: " + c.Logo + "\n" +
		"\tWazuhVersion: " + c.WazuhVersion + "\n" +
		"\tWazuhIP: " + c.WazuhIP + "\n" +
		"\tWazuhPort: " + c.WazuhPort + "\n" +
		"\tWazuhUsername: " + c.WazuhUsername + "\n" +
		"\tIndexerIP: " + c.IndexerIP + "\n" +
		"\tIndexerPort: " + c.IndexerPort + "\n" +
		"}"
}

func CreateClient(db *gorm.DB,
	Name string,
	Logo string,
	WazuhIP string,
	WazuhPort string,
	WazuhUsername string,
	WazuhPassword string,
	IndexerIP string,
	IndexerPort string,
	IndexerUsername string,
	IndexerPassword string,
) (*Client, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	if Name == "" {
		return nil, errors.New("client name cannot be empty")
	}

	if utils.IsValidIpOrDomain(WazuhIP) == false {
		return nil, errors.New("invalid Wazuh IP or domain")
	}

	if utils.IsValidPort(WazuhPort) == false {
		return nil, errors.New("invalid Wazuh port")
	}

	if utils.IsValidIpOrDomain(IndexerIP) == false {
		return nil, errors.New("invalid indexer IP or domain")
	}

	if utils.IsValidPort(IndexerPort) == false {
		return nil, errors.New("invalid indexer port")
	}

	if WazuhUsername == "" || WazuhPassword == "" {
		return nil, errors.New("Wazuh username and password cannot be empty")
	}

	if IndexerUsername == "" || IndexerPassword == "" {
		return nil, errors.New("Indexer username and password cannot be empty")
	}

	client := Client{
		ID:              uuid.String(),
		Name:            Name,
		Logo:            Logo,
		WazuhIP:         WazuhIP,
		WazuhPort:       WazuhPort,
		WazuhUsername:   WazuhUsername,
		WazuhPassword:   WazuhPassword,
		IndexerIP:       IndexerIP,
		IndexerPort:     IndexerPort,
		IndexerUsername: IndexerUsername,
		IndexerPassword: IndexerPassword,
	}
	result := db.Create(&client)

	if result.Error != nil {
		return nil, result.Error
	}

	return &client, nil
}

func GetAllClients(db *gorm.DB) ([]Client, error) {
	var clients []Client
	result := db.Find(&clients)
	if result.Error != nil {
		return nil, result.Error
	}
	return clients, nil
}

func GetClientByID(db *gorm.DB, id string) (*Client, error) {
	client := Client{}
	result := db.First(&client, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &client, nil
}

func GetClientByName(db *gorm.DB, name string) (*Client, error) {
	client := Client{}
	result := db.First(&client, "name = ?", name)
	if result.Error != nil {
		return nil, result.Error
	}
	return &client, nil
}

func DeleteClient(db *gorm.DB, id string) error {
	_, err := GetClientByID(db, id)
	if err != nil {
		return err
	}
	result := db.Delete(&Client{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func EditClient(db *gorm.DB, id string,
	name string,
	logo string,
	wazuhIP string,
	wazuhPort string,
	wazuhUsername string,
	wazuhPassword string,
	indexerIP string,
	indexerPort string,
	indexerUsername string,
	indexerPassword string,
) (Client, error) {
	client := Client{}
	result := db.First(&client, "id = ?", id)
	if result.Error != nil {
		return Client{}, result.Error
	}

	client.Name = name
	client.Logo = logo
	client.WazuhIP = wazuhIP
	client.WazuhPort = wazuhPort
	client.WazuhUsername = wazuhUsername
	client.WazuhPassword = wazuhPassword
	client.IndexerIP = indexerIP
	client.IndexerPort = indexerPort
	client.IndexerUsername = indexerUsername
	client.IndexerPassword = indexerPassword

	result = db.Save(&client)
	if result.Error != nil {
		return Client{}, result.Error
	}
	return client, nil
}
