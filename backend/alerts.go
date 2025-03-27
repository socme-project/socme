package main

import (
	"backend/model"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func (b *Backend) AlertRoutes() {
	// GET /alerts - List all alerts
	b.Router.GET("/alerts", b.userMiddleware, func(c *gin.Context) {
		perPage, _ := strconv.Atoi(c.Query("perPage"))
		page, _ := strconv.Atoi(c.Query("page"))
		severity, _ := c.GetQuery("severity")
		search, _ := c.GetQuery("search")
		filter := Filter{
			Severity: strings.Split(severity, ","),
		}

		alerts, totalNumberOfPages, err := b.SearchAlert(search, filter, perPage, page)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to retrieve alerts.", "error": err.Error()},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{"alerts": alerts, "maxPage": totalNumberOfPages, "message": "Alerts retrieved."},
		)
	})

	// GET /alerts/stats/:severity - Get stats for alerts
	// TODO:  send critical more often, then high a bit less, etc
	b.Router.GET("/alerts/stats/:severity", b.userMiddleware, func(c *gin.Context) {
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
			gin.H{"events": alertsPerHour12, "message": "Stats retrieved."},
		)
	})

	// GET /alerts/:id - Get an alert by ID
	b.Router.GET("/alerts/:id", b.userMiddleware, func(c *gin.Context) {
		id := c.Param("id")
		alert := model.Alert{}

		if err := b.Db.First(&alert, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Alert not found."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"alert": alert, "message": "Alert retrieved."})
	})
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
	// 	query = query.Where("client_id IN ?", filter.ClientID)
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
