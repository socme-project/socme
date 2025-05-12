package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Alert struct {
	ID       uint `gorm:"primaryKey"`
	ClientID uint `                  json:"client_id"`

	ClientName      string    `gorm:"index" json:"client_name"`
	WazuhAlertID    string    `             json:"wazuh_alert_id"`
	RuleID          string    `             json:"rule_id"`
	RuleLevel       uint      `             json:"rule_level"`
	RuleDescription string    `             json:"rule_description"`
	Timestamp       time.Time `             json:"timestamp"`
	RawJSON         string    `             json:"raw_json"`
	Tags            string    `             json:"tags"`
	Sort            int       `             json:"sort"`

	Client *Client `gorm:"constraint:OnDelete:CASCADE;" json:"client,omitempty"`
}

func (c *Client) NewAlert(
	db *gorm.DB,
	wazuhAlertID, ruleID, ruleDescription, rawJSON string,
	sort int,
	timestamp time.Time,
	ruleLevel uint,
) error {
	alert := Alert{
		ClientID:        c.ID,
		ClientName:      c.Name,
		WazuhAlertID:    wazuhAlertID,
		RuleID:          ruleID,
		RuleLevel:       ruleLevel,
		RuleDescription: ruleDescription,
		Timestamp:       timestamp,
		RawJSON:         rawJSON,
		Sort:            sort,
	}

	// TODO: Assign tags here

	result := db.Create(&alert)
	if result.Error != nil {
		return fmt.Errorf("Error while creating the alert.")
	}

	return nil
}
