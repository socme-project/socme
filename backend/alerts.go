package main

import (
	"backend/model"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
	wazuhapi "github.com/socme-project/wazuh-go"
)

func (b *Backend) AlertRoutes() {
	b.Router.GET("/alerts/all", b.userMiddleware, func(c *gin.Context) {
		clientIDStr := c.Query("clientID") // Get ClientID from query params

		var alerts []model.Alert
		query := b.Db.Order("timestamp DESC")

		if clientIDStr != "" {
			clientID, err := strconv.ParseUint(clientIDStr, 10, 32)
			if err == nil {
				query = query.Where("client_id = ?", uint(clientID))
			}
		}

		query.Find(&alerts)

		c.JSON(http.StatusOK, gin.H{"alerts": alerts, "message": "All alerts retrieved"})
	})

	b.Router.GET("/alerts/page", b.userMiddleware, func(c *gin.Context) {
		perPage, _ := strconv.Atoi(c.Query("perPage"))
		page, _ := strconv.Atoi(c.Query("page"))
		severity, _ := c.GetQuery("severity")
		search, _ := c.GetQuery("search")
		filter := Filter{
			Severity: strings.Split(severity, ","),
		}
		var totalNumberOfPages int64 = 0
		b.Db.Model(&model.Alert{}).Count(&totalNumberOfPages)
		totalNumberOfPages = totalNumberOfPages/int64(perPage) + 1

		alerts, totalNumberOfPages, err := b.SearchAlert(search, filter, perPage, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve alerts"})
		}

		c.JSON(
			http.StatusOK,
			gin.H{"alerts": alerts, "maxPage": totalNumberOfPages, "message": "Page retrieved"},
		)
	})

	// TODO: Filter by clientname, if clientname == all, then no filter
	b.Router.GET("/alerts/getlastfive", b.userMiddleware, func(c *gin.Context) {
		alerts := []model.Alert{}
		b.Db.Order("timestamp DESC").Where("rule_level >= ?", 12).Limit(5).Find(&alerts)
		c.JSON(http.StatusOK, gin.H{"alerts": alerts, "message": "Last five alerts retrieved"})
	})

	// TODO: Filter by clientname, if clientname == all, then no filter, send critical more often, then high a bit less, etc
	b.Router.GET("/alerts/last24h/:severity", b.userMiddleware, func(c *gin.Context) {
		severity := c.Param("severity")
		query := b.Db.Model(&model.Alert{})
		switch severity {
		case "low":
			query = query.Where("rule_level <= ?", 6)
		case "medium":
			query = query.Where("rule_level >= ? AND rule_level <= ?", 7, 11)
		case "high":
			query = query.Where("rule_level >= ? AND rule_level <= ?", 12, 14)
		case "critical":
			query = query.Where("rule_level >= ?", 15)
		}

		end := time.Now()
		start := end.Add(-25 * time.Hour)
		b.Logger.Info("Start: ", start, " Now: ", end) // Get the interval

		var alerts []model.Alert
		query.Where("timestamp BETWEEN ? AND ?", start, end).Find(&alerts)
		alertsPerHour := make([]int, 24)

		for _, alert := range alerts {
			hoursAgo := int(end.Sub(alert.Timestamp).Hours())
			if hoursAgo >= 0 && hoursAgo < 24 {
				alertsPerHour[23-hoursAgo]++ // Inverted the index to have the last hour at the end
			}
		}

		var alertsPerHour12 []int

		for i := 0; i < 24; i += 2 {
			alertsPerHour12 = append(alertsPerHour12, alertsPerHour[i]+alertsPerHour[i+1])
		}

		c.JSON(
			http.StatusOK,
			gin.H{"events": alertsPerHour12, "message": "Last 24h alerts retrieved"},
		)
	})
	b.Router.GET("/alerts/id/:id", b.userMiddleware, func(c *gin.Context) {
		id := c.Param("id")
		alert := model.Alert{}
		b.Db.First(&alert, id)
		c.JSON(http.StatusOK, gin.H{"alert": alert, "message": "Alert retrieved"})
	})
}

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
		// Change that
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

type Filter struct {
	Severity   []string
	RuleID     []string
	ClientName []string
	Tags       string
}

// dorks -> ruleid rulelevel description(no need cause already implied)
// search by rule description first, cf https://github.com/lithammer/fuzzysearch
func (b Backend) SearchAlert(
	search string,
	filter Filter,
	perPage, page int,
) ([]model.Alert, int64, error) {
	var alerts []model.Alert
	query := b.Db.Model(&model.Alert{}).Order("timestamp DESC")

	// Filter by severity as string
	if len(filter.Severity) > 0 && filter.Severity[0] != "" {
		severityQuery := b.Db.Where("1 = 0") // Start with a false condition
		for _, severity := range filter.Severity {
			switch severity {
			case "low":
				severityQuery = severityQuery.Or("rule_level <= ?", 6)
			case "medium":
				severityQuery = severityQuery.Or("rule_level >= ? AND rule_level <= ?", 7, 11)
			case "high":
				severityQuery = severityQuery.Or("rule_level >= ? AND rule_level <= ?", 12, 14)
			case "critical":
				severityQuery = severityQuery.Or("rule_level >= ?", 15)
			}
		}
		query = query.Where(severityQuery)
	}

	// Add clientID filtering if you've added it to the Filter struct
	// if filter.ClientID != nil && len(filter.ClientID) > 0 {
	//     query = query.Where("client_id IN ?", filter.ClientID)
	// }

	var totalNumberOfPages int64 = 0
	if search == "" {
		query.Count(&totalNumberOfPages)
		totalNumberOfPages = totalNumberOfPages/int64(perPage) + 1
		if err := query.Limit(perPage).Offset((page - 1) * perPage).Find(&alerts).Error; err != nil {
			return nil, 0, err
		}
	} else {
		rows, err := query.Rows()
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()

		type AlertRank struct {
			alert model.Alert
			rank  int
		}

		alertsRank := []AlertRank{}
		for rows.Next() {
			var alert model.Alert
			err := b.Db.ScanRows(rows, &alert)
			if err != nil {
				return nil, 0, err
			}

			rank := fuzzy.RankMatchNormalizedFold(search, alert.RuleDescription)
			if rank >= 0 {
				alertsRank = append(alertsRank, AlertRank{alert, rank})
			}
		}

		// Sort by rank then put in alerts
		if len(alertsRank) > 0 {
			sort.Slice(alertsRank[:], func(i, j int) bool {
				return alertsRank[i].rank < alertsRank[j].rank
			})

			for _, alertRank := range alertsRank {
				alerts = append(alerts, alertRank.alert)
			}

			totalNumberOfPages = int64(len(alerts))/int64(perPage) + 1

			// Ensure pagination doesn't exceed array bounds
			start := (page - 1) * perPage
			end := page * perPage

			if start < len(alerts) {
				if end > len(alerts) {
					end = len(alerts)
				}
				alerts = alerts[start:end]
			} else {
				alerts = []model.Alert{}
			}
		}
	}

	return alerts, totalNumberOfPages, nil
}
