package main

import (
	"backend/model"
	"log"
	"strconv"
	"time"

	wazuhapi "github.com/socme-project/wazuh-go"
	"gorm.io/gorm"
)

func (b Backend) UpdateAlerts() {
	go func() {
		lastID, err := b.GetLastAlertIdFromDb()
		if err != nil {
			log.Println("Failed to retrieve last alert ID from db:", err)
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
	}()
}

func (b Backend) AddAlertToDb(alerts []wazuhapi.Alert) error {
	for _, alert := range alerts {
		timestamp, err := time.Parse(time.RFC3339, alert.Timestamp)
		if err != nil {
			return err
		}
		alert.Timestamp = timestamp.Format(time.RFC3339)
		// TODO: check if we need to RFC3339, or if there is another universal default format

		err = model.NewAlert(
			b.Db,
			alert.WazuhAlertID,
			alert.RuleID,
			alert.RuleDescription,
			alert.Timestamp,
			alert.RawJSON,
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
	return strconv.Atoi(alert.WazuhAlertID)
}

type Filter struct {
	Severity   []string
	RuleID     []string
	ClientName []string
	Tags       []string
}

// dorks -> ruleid rulelevel description(no need cause already implied)
// search by rule description first, cf https://github.com/lithammer/fuzzysearch
func (b Backend) SearchAlert(input string, filter Filter) ([]model.Alert, error) {
	var alerts []model.Alert
	query := b.Db.Model(&model.Alert{})

	// if input != "" { // only for fzf
	// 	query = query.Where(
	// 		"rule_description ILIKE ? OR raw_json ILIKE ?",
	// 		"%"+input+"%",
	// 		"%"+input+"%",
	// 	)
	// }

	// Filter by severity as string
	if len(filter.Severity) > 0 {
		query = query.Where(func(db *gorm.DB) *gorm.DB {
			for _, severity := range filter.Severity {
				switch severity {
				case "low":
					db = db.Or("rule_level <= ?", 6)
				case "medium":
					db = db.Or("rule_level >= ? AND rule_level <= ?", 7, 11)
				case "high":
					db = db.Or("rule_level >= ? AND rule_level <= ?", 12, 14)
				case "critical":
					db = db.Or("rule_level >= ?", 15)
				}
			}
			return db
		})
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

	if err := query.Find(&alerts).Error; err != nil {
		return nil, err
	}

	return alerts, nil
}

// func (b Backend) SearchPaginated() {
