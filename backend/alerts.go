package main

import (
	"backend/model"
	"log"
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
		alerts := []model.Alert{}
		b.Db.Order("timestamp DESC").Find(&alerts)
		c.JSON(http.StatusOK, gin.H{"alerts": alerts, "message": "All alerts retrieved"})
	})

	b.Router.GET("/alerts/page", b.userMiddleware, func(c *gin.Context) {
		perPage, _ := strconv.Atoi(c.Query("perPage"))
		page, _ := strconv.Atoi(c.Query("page"))
		severity, _ := c.GetQuery("severity")
		search, _ := c.GetQuery("search")
		filter := Filter{
			Severity: strings.Split(severity, ","), // split by comma

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

	b.Router.GET("/alerts/getlastfive", b.userMiddleware, func(c *gin.Context) {
		alerts := []model.Alert{}
		b.Db.Order("timestamp DESC").Where("rule_level >= ?", 12).Limit(5).Find(&alerts)
		c.JSON(http.StatusOK, gin.H{"alerts": alerts, "message": "Last five alerts retrieved"})
	})

	b.Router.GET("/alerts/last24h/:severity", b.userMiddleware, func(c *gin.Context) {
		severity := c.Param("severity")
		alerts := []model.Alert{}
		query := b.Db.Model(&model.Alert{})
		switch severity {
		case "low":
			query = query.Or("rule_level <= ?", 6)
		case "medium":
			query = query.Or("rule_level >= ? AND rule_level <= ?", 7, 11)
		case "high":
			query = query.Or("rule_level >= ? AND rule_level <= ?", 12, 14)
		case "critical":
			query = query.Or("rule_level >= ?", 15)
		}
		// array of 12 length, 0 is the oldest... to display in dashboard, medium has its own bars, high has its own bars, etc
		for i := 0; i < 12; i++ {
			var alertsInHour []model.Alert
			query.Where("timestamp BETWEEN ? AND ?", time.Now().Add(-time.Hour*time.Duration(i+1)), time.Now().Add(-time.Hour*time.Duration(i))).
				Find(&alertsInHour)
			alerts = append(alerts, alertsInHour...)
		}
		c.JSON(http.StatusOK, gin.H{"alerts": alerts, "message": "Last 24h alerts retrieved"})
	})
	b.Router.GET("/alerts/id/:id", b.userMiddleware, func(c *gin.Context) {
		id := c.Param("id")
		alert := model.Alert{}
		b.Db.First(&alert, id)
		c.JSON(http.StatusOK, gin.H{"alert": alert, "message": "Alert retrieved"})
	})
}

func (b Backend) UpdateAlerts() {
	lastID, err := b.GetLastAlertIdFromDb()
	if err != nil && err.Error() == "record not found" {
		lastID = 0
	} else if err != nil {
		log.Println("Failed to retrieve last alert ID from db:", err)
	}

	err = b.Wazuh.RefreshToken()
	if err != nil {
		log.Fatal("Failed to refresh token:", err)
	}

	for {
		alerts, newLastID, err := b.Wazuh.GetAlerts(lastID)
		if err != nil {
			log.Println("Failed to retrieve alerts:", err)
		}
		if len(alerts) > 0 {
			err := b.AddAlertToDb(alerts)
			if err != nil {
				log.Println("Failed to add alerts to db:", err)
			}
		}
		lastID = newLastID
		time.Sleep(b.AlertRetrievalInterval)
	}
}

func (b Backend) AddAlertToDb(alerts []wazuhapi.Alert) error {
	layout := "2006-01-02T15:04:05.000-0700"
	for _, alert := range alerts {
		timestamp, err := time.Parse(layout, alert.Timestamp)
		if err != nil {
			return err
		}
		// TODO: check if we need to RFC3339, or if there is another universal default format

		err = model.NewAlert(
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

func (b Backend) GetLastAlertIdFromDb() (lastID int, err error) {
	var alert model.Alert
	result := b.Db.Order("timestamp DESC").First(&alert) // check if result is logic
	if result.Error != nil {
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
	if len(filter.Severity) > 0 {
		for _, severity := range filter.Severity {
			switch severity {
			case "low":
				query = query.Or("rule_level <= ?", 6)
			case "medium":
				query = query.Or("rule_level >= ? AND rule_level <= ?", 7, 11)
			case "high":
				query = query.Or("rule_level >= ? AND rule_level <= ?", 12, 14)
			case "critical":
				query = query.Or("rule_level >= ?", 15)
			}
		}
	}

	// if len(filter.RuleID) > 0 {
	// 	query = query.Where("rule_id IN ?", filter.RuleID)
	// }
	//
	// if len(filter.ClientName) > 0 {
	// 	query = query.Where("client_name IN ?", filter.ClientName)
	// }
	//
	// if len(filter.Tags) > 0 {
	// 	for _, tag := range filter.Tags {
	// 		query = query.Where("tags ILIKE ?", "%"+tag+"%")
	// 	}
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
			_ = b.Db.ScanRows(rows, &alert)
			rank := fuzzy.RankMatchNormalizedFold(search, alert.RuleDescription)
			if rank >= 0 {
				alertsRank = append(alertsRank, AlertRank{alert, rank})
			}
		}
		// sort by rank then put in alerts
		if len(alertsRank) > 0 {
			sort.Slice(alertsRank[:], func(i, j int) bool {
				return alertsRank[i].rank < alertsRank[j].rank
			})
		} else {
			return []model.Alert{}, 0, nil
		}
		for _, alertRank := range alertsRank {
			alerts = append(alerts, alertRank.alert)
		}
		totalNumberOfPages = int64(len(alerts))/int64(perPage) + 1
		if len(alerts) > perPage {
			alerts = alerts[(page-1)*perPage : page*perPage]
		}
	}

	return alerts, totalNumberOfPages, nil
}

// func (b Backend) SearchPaginated() {
