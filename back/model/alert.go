package model

import (
	"strconv"
	"time"
)

type Alert struct {
	ID              uint      `gorm:"primaryKey"`
	ClientID        string    `gorm:"not null"` // ClientID is the ID of the client that this alert belongs to
	WazuhAlertID    string    `             json:"wazuh_alert_id"`
	RuleID          string    `             json:"rule_id"`
	RuleLevel       uint      `             json:"rule_level"`
	RuleDescription string    `             json:"rule_description"`
	Timestamp       time.Time `             json:"timestamp"`
	Tags            string    `             json:"tags"`     // Tags is category of alerts
	Sort            int       `             json:"sort"`     // Sort is an index
	RawJSON         string    `             json:"raw_json"` // Raw alert from Wazuh
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

// NewAlert()
