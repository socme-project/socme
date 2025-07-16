package routes

import (
	"strconv"
	"time"

	wazuhapi "github.com/socme-project/wazuh-go"
	"github.com/socme-projects/backend/model"
)

func (r routerType) UpdateAlertsForClient(client model.Client) error {
	r.Logger.Info("-- Retrieving alerts for " + client.Name)
	lastID, err := r.GetLastAlertIdFromDb(client.ID)

	if err != nil && err.Error() == "record not found" {
		lastID = 0
	} else if err != nil {
		r.Logger.Error("Failed to retrieve last alert ID from db: " + err.Error())
		return err
	}
	r.Logger.Info("Last ID: " + strconv.Itoa(lastID))

	wazuhClient := wazuhapi.WazuhAPI{
		Host:     client.Host,
		Port:     client.WazuhPort,
		Username: client.WazuhUsername,
		Password: client.WazuhPassword,
		Indexer: wazuhapi.Indexer{
			Host:     client.Host,
			Port:     client.IndexerPort,
			Username: client.IndexerUsername,
			Password: client.IndexerPassword,
		},
		Insecure: true,
	}

	err = wazuhClient.RefreshToken()
	if err != nil {
		r.Logger.Error("Failed to refresh token: ", err.Error())
		return err
	}

	alerts, _, err := wazuhClient.GetAlerts(lastID)
	if err != nil {
		r.Logger.Error("Failed to retrieve alerts: " + err.Error())
		return err
	} else if len(alerts) == 0 {
		r.Logger.Info("No new alerts found for client " + client.Name)
		return nil
	}

	err = r.AddAlertToDb(alerts, client.ID)
	if err != nil {
		r.Logger.Error("Failed to add alerts to db:", err)
		return err
	}

	lastALert, err := time.Parse("2006-01-02T15:04:05.000-0700", alerts[len(alerts)-1].Timestamp)
	if err != nil {
		r.Logger.Error("Failed to parse last alert timestamp: " + err.Error())
		return err
	}
	err = model.EditLastAlert(r.Db, client, lastALert)

	if err != nil {
		r.Logger.Error("Failed to update last alert timestamp in db: " + err.Error())
		return err
	}

	// TODO: fastetch here (waiting for opsme integration)
	version, err := wazuhClient.GetApiVersion()
	if err != nil {
		r.Logger.Error("Failed to retrieve API version: ", err.Error())
		return err
	}
	err = model.EditClientVersion(r.Db, client.ID, version)
	if err != nil {
		r.Logger.Error("Failed to update client version in db: ", err.Error())
		return err
	}

	agents, err := wazuhClient.GetAgents()
	if err != nil {
		r.Logger.Error("Failed to retrieve agents status: ", err.Error())
		return err
	}
	err = model.EditClientAgents(r.Db, client.ID, agents.Active, agents.Disconnected)
	if err != nil {
		r.Logger.Error("Failed to update client agents in db: ", err.Error())
		return err
	}

	return nil
}

func (r routerType) UpdateAlerts() {
	r.Logger.Info("Starting alert retrieval")
	for {
		clients, err := model.GetAllClients(r.Db)
		if err != nil {
			r.Logger.Error("Failed to retrieve clients from db: ", err)
			time.Sleep(r.RefreshRate)
			continue
		}
		r.Logger.Info("Retrieving alerts for", len(clients), "clients: ", clients)
		for _, client := range clients {
			go func() {
				err := r.UpdateAlertsForClient(client)
				if err != nil {
					model.EditClientStatus(r.Db, client.ID, false)
				}
				model.EditClientStatus(r.Db, client.ID, true)
			}()
		}
		time.Sleep(r.RefreshRate)
	}
}

func (r routerType) AddAlertToDb(alerts []wazuhapi.Alert, clientID string) error {
	client, err := model.GetClientByID(r.Db, clientID)
	if err != nil {
		return err
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
			timestamp.UTC(),
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
