package routes

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/socme-projects/backend/model"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Filter struct {
	Severity []string
	RuleID   []string
	ClientID string
	Tags     string
}

func (r *routerType) alertsRoutes() {
	g := r.R.Group("/alerts")

	g.GET("", r.RoleMiddleware(), func(c *gin.Context) {
		perPageStr := c.DefaultQuery("perPage", "10")
		perPage, err := strconv.Atoi(perPageStr)
		if err != nil || perPage <= 0 {
			perPage = 10
		}

		page, _ := strconv.Atoi(c.Query("page"))
		severity, _ := c.GetQuery("severity")
		search, _ := c.GetQuery("search")
		client, _ := c.GetQuery("client")
		preload := false
		if p, _ := c.GetQuery("preload"); p == "true" {
			preload = true
		}

		filter := Filter{
			Severity: strings.Split(severity, ","),
			ClientID: client,
		}

		alerts, totalNumberOfPages, err := r.SearchAlert(search, filter, perPage, page, preload)
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

	g.GET(":id", r.RoleMiddleware(), func(c *gin.Context) {
		id := c.Param("id")
		idUint, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"message": "Invalid alert ID format.", "error": err.Error()},
			)
			return
		}
		alert, err := model.GetAlertByID(r.Db, uint(idUint))
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to retrieve alert.", "error": err.Error()},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{"alert": alert, "message": "Alert retrieved."},
		)
	})

	g.GET("stats/:severity", r.RoleMiddleware(), func(c *gin.Context) {
		severity := c.Param("severity")
		query := r.Db.Model(&model.Alert{})
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

		end := time.Now().UTC()
		start := end.Add(-24 * time.Hour).UTC()

		var alerts []model.Alert
		query.Where("timestamp BETWEEN ? AND ?", start, end).Find(&alerts)
		// FIXME: Is the timestamp thing good ?

		// fmt.Println("\n\n---------------------------------------------")
		// fmt.Println(len(alerts))

		alertsPerHour := make([]int, 24)

		for _, alert := range alerts {
			hoursAgo := int(end.Sub(alert.Timestamp).Hours())
			if hoursAgo >= 0 && hoursAgo < 24 {
				alertsPerHour[23-hoursAgo]++ // Inverted the index to have the last hour at the end
			}
		}

		var alertsPerHour12 []int

		for i := 0; i < 24; i += 4 {
			alertsPerHour12 = append(alertsPerHour12, alertsPerHour[i]+alertsPerHour[i+1]+alertsPerHour[i+2]+alertsPerHour[i+3])
		}

		c.JSON(
			http.StatusOK,
			gin.H{"events": alertsPerHour12, "message": "Stats retrieved."},
		)
	})

}

// dorks -> ruleid rulelevel description(no need cause already implied)
// search by rule description first, cf https://github.com/lithammer/fuzzysearch
func (r routerType) SearchAlert(
	search string,
	filter Filter,
	perPage, page int,
	preload bool,
) ([]model.Alert, int64, error) {
	db := r.Db
	if preload {
		db = db.Preload("Client")
	}

	var alerts []model.Alert
	query := db.Model(&model.Alert{}).Order("timestamp DESC")

	// Filter by severity as string
	if len(filter.Severity) > 0 && filter.Severity[0] != "" {
		severityQuery := db.Where("1 = 0") // Start with a false condition
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

	if filter.ClientID != "" {
		query.Where("client_id = ?", filter.ClientID)
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
			err := db.ScanRows(rows, &alert)
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
			sort.Slice(alertsRank, func(i, j int) bool {
				if alertsRank[i].rank == alertsRank[j].rank {
					return alertsRank[i].alert.Timestamp.After(alertsRank[j].alert.Timestamp)
				}
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
