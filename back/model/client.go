package model

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/socme-project/backend/utils"
	"gorm.io/gorm"
)

type Client struct {
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
	Logo string

	// Settings
	Host string

	SshPort     string
	SshUsername string
	SshPassword string

	WazuhPort     string
	WazuhUsername string
	WazuhPassword string

	IndexerPort     string
	IndexerUsername string
	IndexerPassword string

	// Auto updated
	LastAlert          time.Time
	ConnectedAgents    int
	DisconnectedAgents int
	Alerts             []Alert `gorm:"foreignKey:ClientID"`
	WazuhIsAlive       bool
	WazuhVersion       string

	Information struct {
		Os     string
		Host   string
		Kernel string
		CPU    string
		GPU    string

		IP     string
		Uptime string

		Disk     string
		Memory   string
		Swap     string
		CPUUsage string
	}
}

func (c Client) String() string {
	return "Client{\n" +
		"\tID: " + c.ID + "\n" +
		"\tName: " + c.Name + "\n" +
		"\tLogo: " + c.Logo + "\n" +
		"\tWazuhVersion: " + c.WazuhVersion + "\n" +
		"\tWazuhPort: " + c.WazuhPort + "\n" +
		"\tWazuhUsername: " + c.WazuhUsername + "\n" +
		"\tIndexerPort: " + c.IndexerPort + "\n" +
		"\tIndexerUsername: " + c.IndexerUsername + "\n" +
		"\tLastAlert: " + c.LastAlert.String() + "\n" +
		"\tConnectedAgents: " + strconv.Itoa(c.ConnectedAgents) + "\n" +
		"\tDisconnectedAgents: " + strconv.Itoa(c.DisconnectedAgents) + "\n" +
		"\tWazuhIsAlive: " + strconv.FormatBool(c.WazuhIsAlive) + "\n" +
		"\tHostIP: " + c.Host + "\n" +
		"\tSshPort: " + c.SshPort + "\n" +
		"\tSshUsername: " + c.SshUsername + "\n" +
		"\tSshPassword: " + c.SshPassword + "\n" +
		"}"
}

func CreateClient(db *gorm.DB,
	Name string,
	Logo string,

	Host string,

	SshPort string,
	SshUsername string,
	SshPassword string,

	WazuhPort string,
	WazuhUsername string,
	WazuhPassword string,

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

	if utils.IsValidIpOrDomain(Host) == false {
		return nil, errors.New("invalid host IP or domain")
	}

	if utils.IsValidPort(WazuhPort) == false {
		return nil, errors.New("invalid Wazuh port")
	}

	if utils.IsValidPort(SshPort) == false {
		return nil, errors.New("invalid SSH port")
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

	if SshUsername == "" || SshPassword == "" {
		return nil, errors.New("SSH username and password cannot be empty")
	}

	client := Client{
		ID:   uuid.String(),
		Name: Name,
		Logo: Logo,

		Host: Host,

		SshPort:     SshPort,
		SshUsername: SshUsername,
		SshPassword: SshPassword,

		WazuhPort:     WazuhPort,
		WazuhUsername: WazuhUsername,
		WazuhPassword: WazuhPassword,

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

func EditLastAlert(db *gorm.DB, client Client, lastAlert time.Time) error {
	client.LastAlert = lastAlert
	result := db.Save(&client)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func EditClient(db *gorm.DB, id string,
	Name string,
	Logo string,

	Host string,

	SshPort string,
	SshUsername string,
	SshPassword string,

	WazuhPort string,
	WazuhUsername string,
	WazuhPassword string,

	IndexerPort string,
	IndexerUsername string,
	IndexerPassword string,
) (Client, error) {
	client := Client{}
	result := db.First(&client, "id = ?", id)
	if result.Error != nil {
		return Client{}, result.Error
	}

	client.Name = Name
	client.Logo = Logo

	client.Host = Host

	client.SshPort = SshPort
	client.SshUsername = SshUsername
	client.SshPassword = SshPassword

	client.WazuhPort = WazuhPort
	client.WazuhUsername = WazuhUsername
	client.WazuhPassword = WazuhPassword

	client.IndexerPort = IndexerPort
	client.IndexerUsername = IndexerUsername
	client.IndexerPassword = IndexerPassword

	result = db.Save(&client)
	if result.Error != nil {
		return Client{}, result.Error
	}
	return client, nil
}

func EditClientAgents(db *gorm.DB, id string, connectedAgents int, disconnectedAgents int) error {
	client := Client{}
	result := db.First(&client, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	client.ConnectedAgents = connectedAgents
	client.DisconnectedAgents = disconnectedAgents
	result = db.Save(&client)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func EditClientVersion(db *gorm.DB, id string, version string) error {
	client := Client{}
	result := db.First(&client, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	client.WazuhVersion = version
	result = db.Save(&client)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func EditClientStatus(db *gorm.DB, id string, isAlive bool) error {
	client := Client{}
	result := db.First(&client, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	client.WazuhIsAlive = isAlive
	result = db.Save(&client)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func EditClientInformation(db *gorm.DB, id string,
	Os string,
	Host string,
	Kernel string,
	CPU string,
	GPU string,

	IP string,
	Uptime string,

	Disk string,
	Memory string,
	Swap string,
	CPUUsage string,
) error {
	client := Client{}
	result := db.First(&client, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	client.Information.Os = Os
	client.Information.Host = Host
	client.Information.Kernel = Kernel
	client.Information.CPU = CPU
	client.Information.GPU = GPU

	client.Information.IP = IP
	client.Information.Uptime = Uptime

	client.Information.Disk = Disk
	client.Information.Memory = Memory
	client.Information.Swap = Swap
	client.Information.CPUUsage = CPUUsage

	result = db.Save(&client)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
