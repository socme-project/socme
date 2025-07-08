package model

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Alerts []Alert

type Alert struct {
	ID              uint `gorm:"primaryKey"`
	WazuhAlertID    string
	RuleID          string
	RuleLevel       uint
	RuleDescription string
	Timestamp       time.Time
	Tags            string // Tags is category of alerts
	Sort            int    // Sort is an index
	RawJSON         string // Raw alert from Wazuh
	ClientID        string
	Client          Client `gorm:"foreignKey:ClientID"`
}

func (a Alert) String() string {
	return "Alert{\n" +
		"\tID: " + strconv.Itoa(int(a.ID)) + "\n" +
		"\tClientID: " + a.ClientID + "\n" +
		"\tWazuhAlertID: " + a.WazuhAlertID + "\n" +
		"\tRuleID: " + a.RuleID + "\n" +
		"\tRuleLevel: " + strconv.Itoa(int(a.RuleLevel)) + "\n" +
		"\tRuleDescription: " + a.RuleDescription + "\n" +
		"\tTimestamp: " + a.Timestamp.String() + "\n" +
		"\tTags: " + a.Tags + "\n" +
		"\tSort: " + strconv.Itoa(a.Sort) + "\n" +
		"\tRawJSON: " + a.RawJSON + "\n" +
		"}"
}

func CreateAlert(
	db *gorm.DB, clientid string, wazuhAlertID, ruleID, ruleDescription, rawJSON string,
	ruleLevel uint, timestamp time.Time, sort int,
) (*Alert, error) {
	alert := Alert{
		WazuhAlertID:    wazuhAlertID,
		RuleID:          ruleID,
		RuleLevel:       ruleLevel,
		RuleDescription: ruleDescription,
		Timestamp:       timestamp,
		RawJSON:         rawJSON,
		ClientID:        clientid,
		Sort:            sort,
	}

	db.Create(&alert)

	return &alert, nil
}

func GetAlertsByClientID(db *gorm.DB, clientID string) ([]Alert, error) {
	var alerts []Alert
	result := db.Where("client_id = ?", clientID).Find(&alerts)
	if result.Error != nil {
		return nil, result.Error
	}
	return alerts, nil
}

func GetAllAlerts(db *gorm.DB) ([]Alert, error) {
	var alerts []Alert
	result := db.Find(&alerts)
	if result.Error != nil {
		return nil, result.Error
	}
	return alerts, nil
}

func GetAlertsFromClient(db *gorm.DB, clientID string) (Alerts, error) {
	var alerts []Alert
	result := db.Where("client_id = ?", clientID).Find(&alerts)
	if result.Error != nil {
		return nil, result.Error
	}
	return alerts, nil
}

func GetAlertByID(db *gorm.DB, id uint) (*Alert, error) {
	alert := Alert{}
	result := db.Preload("Client").First(&alert, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &alert, nil
}

func (a Alerts) FilterBySeverity(severity string) []Alert {
	fun := func(alert Alert) bool { return false }

	switch severity {
	case "low":
		fun = func(alert Alert) bool {
			if alert.RuleLevel <= 6 {
				return true
			}
			return false
		}
	case "medium":
		fun = func(alert Alert) bool {
			if alert.RuleLevel >= 7 && alert.RuleLevel <= 11 {
				return true
			}
			return false
		}
	case "high":
		fun = func(alert Alert) bool {
			if alert.RuleLevel >= 12 && alert.RuleLevel <= 14 {
				return true
			}
			return false
		}
	case "critical":
		fun = func(alert Alert) bool {
			if alert.RuleLevel >= 15 {
				return true
			}
			return false
		}
	}

	filteredAlerts := []Alert{}

	for i := range a {
		if fun(a[i]) {
			filteredAlerts = append(filteredAlerts, a[i])
		}
	}

	return filteredAlerts
}
