package main

import (
	"backend/model"
	"fmt"
	"strconv"
	"time"

	wazuhapi "github.com/socme-project/wazuh-go"
)

func (b Backend) UpdateAlertsForClient(client model.Client) {
	b.Logger.Info("-- Retrieving alerts for " + client.Name)
	lastID, err := b.GetLastAlertIdFromDb(client.ID)

	if err != nil && err.Error() == "record not found" {
		lastID = 0
	} else if err != nil {
		b.Logger.Error("Failed to retrieve last alert ID from db: " + err.Error())
		return
	}
	b.Logger.Info("Last ID: " + strconv.Itoa(lastID))

	wazuhClient := wazuhapi.WazuhAPI{
		Host:     client.WazuhIP,
		Port:     client.WazuhPort,
		Username: client.WazuhUsername,
		Password: client.WazuhPassword,
		Indexer: wazuhapi.Indexer{
			Host:     client.IndexerIP,
			Port:     client.IndexerPort,
			Username: client.IndexerUsername,
			Password: client.IndexerPassword,
		},
		Insecure: true,
	}

	if wazuhClient.RefreshToken() != nil {
		b.Logger.Error("Failed to refresh token: " + err.Error())
		return
	}

	alerts, _, err := wazuhClient.GetAlerts(lastID)
	if err != nil {
		b.Logger.Error("Failed to retrieve alerts: " + err.Error())
		return
	} else if len(alerts) == 0 {
		return
	}

	err = b.AddAlertToDb(alerts, client.ID)
	if err != nil {
		b.Logger.Error("Failed to add alerts to db:", err)
		return
	}
}

func (b Backend) UpdateAlerts() {
	b.Logger.Info("Starting alert retrieval")
	for {
		clients := model.GetAllClients(b.Db)
		b.Logger.Info("Retrieving alerts for", len(clients), "clients: ", clients)
		for _, client := range clients {
			go b.UpdateAlertsForClient(client)
		}
		time.Sleep(b.RefreshRate)
	}
}

func (b Backend) AddAlertToDb(alerts []wazuhapi.Alert, clientID uint) error {
	var client model.Client
	if err := b.Db.First(&client, clientID).Error; err != nil {
		return fmt.Errorf("client not found: %w", err)
	}

	layout := "2006-01-02T15:04:05.000-0700"
	b.Logger.Info("Adding alerts for client: ", client.Name)
	for _, alert := range alerts {
		timestamp, err := time.Parse(layout, alert.Timestamp)
		if err != nil {
			return err
		}

		err = client.NewAlert(
			b.Db,
			alert.WazuhAlertID,
			alert.RuleID,
			alert.RuleDescription,
			alert.RawJSON,
			alert.Sort,
			timestamp,
			alert.RuleLevel,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b Backend) GetLastAlertIdFromDb(clientID uint) (int, error) {
	var alert model.Alert
	result := b.Db.Order("timestamp DESC, sort DESC").
		Where("client_id = ?", clientID).
		First(&alert)
	if result.Error != nil {
		b.Logger.Error("Error while getting last alert ID from db: ", result.Error)
		return 0, result.Error
	}
	return alert.Sort, nil
}
