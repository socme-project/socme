package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Alert struct {
	ID         uint   `gorm:"primaryKey"`
	ClientName string `json:"client_name"`

	WazuhAlertID    string    `json:"wazuh_alert_id"`
	RuleID          string    `json:"rule_id"`
	RuleLevel       uint      `json:"rule_level"`
	RuleDescription string    `json:"rule_description"`
	Timestamp       time.Time `json:"timestamp"`
	RawJSON         string    `json:"raw_json"`
	Tags            string    `json:"tags"`
	Sort            int       `json:"sort"`
}

func NewAlert(
	db *gorm.DB,
	wazuhAlertID, ruleID, ruleDescription, rawJSON string,
	sort int,
	timestamp time.Time,
	ruleLevel uint,
	clientName string,
) error {
	alert := Alert{
		// ClientName:      clientName,

		WazuhAlertID:    wazuhAlertID,
		RuleID:          ruleID,
		RuleLevel:       ruleLevel,
		RuleDescription: ruleDescription,
		Timestamp:       timestamp,
		RawJSON:         rawJSON,
		Sort:            sort,
	}

	// Assign tags

	result := db.Create(&alert)
	if result.Error != nil {
		return fmt.Errorf("Error while creating the alert")
	}

	return nil
}

func GetAllAlerts(db *gorm.DB) []Alert {
	var alerts []Alert
	db.Find(&alerts)

	return alerts
}

func GetAlertsByClient(db *gorm.DB, clientName string) ([]Alert, error) {
	var alerts []Alert
	db.Find(&alerts, "client_name = ?", clientName)
	if len(alerts) == 0 {
		return nil, fmt.Errorf("No alerts found for this client")
	}

	return alerts, nil
}

func GetAlertsByRuleLevel(db *gorm.DB, ruleLevel int) ([]Alert, error) {
	var alerts []Alert
	db.Find(&alerts, "rule_level = ?", ruleLevel)
	if len(alerts) == 0 {
		return nil, fmt.Errorf("No alerts found for this rule level")
	}

	return alerts, nil
}

func GetAlertsByFzf(db *gorm.DB, clientName, ruleDescription string) ([]Alert, error) { // TO TEST
	var alerts []Alert
	db.Find(&alerts, "client_name = ? AND rule_description = ?", clientName, ruleDescription)
	if len(alerts) == 0 {
		return nil, fmt.Errorf("No alerts found for this client and rule description")
	}

	return alerts, nil
}

func GetAllPagniatedAlerts(db *gorm.DB, perPage, page int) ([]Alert, int) {
	var alerts []Alert
	db.Limit(perPage).Offset(page * perPage).Find(&alerts)

	var total int64
	db.Model(&Alert{}).Count(&total)

	return alerts, int(total)
}
