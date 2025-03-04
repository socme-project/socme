package model

type Alert struct {
	ID              uint   `gorm:"primaryKey"`
	ClientName      string `json:"client_name"`      // Agent name
	RuleLevel       int    `json:"rule_level"`       // Severity level
	RuleDescription string `json:"rule_description"` // Description of the alert
	Timestamp       string // Event timestamp

	RawJSON string `json:"raw_json"` // Full event data
	//IntegrityLevel  string `json:"integrity_level"`
}

func NewAlert(clientName, ruleDescription, timestamp, rawJSON string, ruleLevel int) Alert {
	return Alert{

		ClientName:      clientName,
		RuleLevel:       ruleLevel,
		RuleDescription: ruleDescription,
		Timestamp:       timestamp,
		RawJSON:         rawJSON,
	}
}

// getAll
// getAlertsByClient
// getAlertsByRuleLevel
// getAlertsByFzf

// getAllPaginated - perPage, page //filtre
// return alerts, nb alertes total
