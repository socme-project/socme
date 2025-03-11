package model

import (
	"fmt"

	"gorm.io/gorm"
)

type Alert struct {
	ID         uint   `gorm:"primaryKey"`
	ClientName string `json:"client_name"`

	WazuhAlertID    string `json:"wazuh_alert_id"`
	RuleID          string `json:"rule_id"`
	RuleLevel       uint   `json:"rule_level"`
	RuleDescription string `json:"rule_description"`
	Timestamp       string `json:"timestamp"`
	RawJSON         string `json:"raw_json"`
}

func NewAlert(
	db *gorm.DB,
	clientName, ruleDescription, timestamp, rawJSON string,
	ruleLevel int,
) error {
	alert := Alert{
		ClientName:      clientName,
		RuleDescription: ruleDescription,
		Timestamp:       timestamp,
		RawJSON:         rawJSON,
	}

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
