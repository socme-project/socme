package routes

import (
	"fmt"
	"strconv"
	"time"

	wazuhapi "github.com/socme-project/wazuh-go"
	"github.com/socme-projects/backend/model"
)

func (r routerType) UpdateAlertsForClient(client model.Client) {
	r.Logger.Info("-- Retrieving alerts for " + client.Name)
	lastID, err := r.GetLastAlertIdFromDb(client.ID)

	if err != nil && err.Error() == "record not found" {
		lastID = 0
	} else if err != nil {
		r.Logger.Error("Failed to retrieve last alert ID from db: " + err.Error())
		return
	}
	r.Logger.Info("Last ID: " + strconv.Itoa(lastID))

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
		r.Logger.Error("Failed to refresh token: " + err.Error())
		return
	}

	alerts, _, err := wazuhClient.GetAlerts(lastID)
	if err != nil {
		r.Logger.Error("Failed to retrieve alerts: " + err.Error())
		return
	} else if len(alerts) == 0 {
		return
	}

	err = r.AddAlertToDb(alerts, client.ID)
	if err != nil {
		r.Logger.Error("Failed to add alerts to db:", err)
		return
	}
}

func (r routerType) UpdateAlerts() {
	r.Logger.Info("Starting alert retrieval")
	for {
		clients, err := model.GetAllClients(r.Db)
		if err != nil {
			r.Logger.Error("Failed to retrieve clients from db: ", err)
			time.Sleep(r.RefreshRate)
		}
		r.Logger.Info("Retrieving alerts for", len(clients), "clients: ", clients)
		for _, client := range clients {
			go r.UpdateAlertsForClient(client)
		}
		time.Sleep(r.RefreshRate)
	}
}

func (r routerType) AddAlertToDb(alerts []wazuhapi.Alert, clientID string) error {
	var client model.Client
	if err := r.Db.First(&client, clientID).Error; err != nil {
		return fmt.Errorf("client not found: %w", err)
	}

	layout := "2006-01-02T15:04:05.000-0700"
	r.Logger.Info("Adding alerts for client: ", client.Name)
	for _, alert := range alerts {
		timestamp, err := time.Parse(layout, alert.Timestamp)
		if err != nil {
			return err
		}

		_, err = model.CreateAlert(
			r.Db,
			clientID,
			alert.WazuhAlertID,
			alert.RuleID,
			alert.RuleDescription,
			alert.RawJSON,
			alert.RuleLevel,
			timestamp,
			alert.Sort,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r routerType) GetLastAlertIdFromDb(clientID string) (int, error) {
	var alert model.Alert
	result := r.Db.Order("timestamp DESC, sort DESC").
		Where("client_id = ?", clientID).
		First(&alert)
	if result.Error != nil {
		r.Logger.Error("Error while getting last alert ID from db: ", result.Error)
		return 0, result.Error
	}
	return alert.Sort, nil
}
