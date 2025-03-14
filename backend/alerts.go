package main

import (
	"backend/model"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
		filter := Filter{
			Severity: strings.Split(severity, ","), // split by comma

		}
		var totalNumberOfPages int64 = 0
		b.Db.Model(&model.Alert{}).Count(&totalNumberOfPages)
		totalNumberOfPages = totalNumberOfPages/int64(perPage) + 1

		alerts, err := b.SearchAlert("", filter, perPage, page)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve alerts"})
		}

		c.JSON(
			http.StatusOK,
			gin.H{"alerts": alerts, "maxPage": totalNumberOfPages, "message": "Page retrieved"},
		)
	})

	b.Router.GET("/alerts/:id", b.userMiddleware, func(c *gin.Context) {
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
	input string,
	filter Filter,
	perPage, page int,
) ([]model.Alert, error) {
	var alerts []model.Alert
	query := b.Db.Model(&model.Alert{}).Order("timestamp DESC")

	// if input != "" { // only for fzf
	// 	query = query.Where(
	// 		"rule_description ILIKE ? OR raw_json ILIKE ?",
	// 		"%"+input+"%",
	// 		"%"+input+"%",
	// 	)
	// }

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

	if err := query.Limit(perPage).Offset((page - 1) * perPage).Find(&alerts).Error; err != nil {
		return nil, err
	}

	return alerts, nil
}

// func (b Backend) SearchPaginated() {
