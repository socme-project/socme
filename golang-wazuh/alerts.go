package wazuhapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Alert struct {
	ID              uint   `json:"id"`
	ClientName      string `json:"client_name"`
	RuleLevel       int    `json:"rule_level"`
	RuleDescription string `json:"rule_description"`
	Timestamp       string `json:"timestamp"`
	RawJSON         string `json:"raw_json"`
}

func (w *WazuhAPI) GetAlerts() ([]Alert, error) {
	resp, err := http.Get(
		"https://10.8.178.20:9200/wazuh-alerts-*/_search/?size=10000",
	)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data struct {
			Items []map[string]interface{} `json:"items"`
		} `json:"data"`
	}

	var response Response
	err = json.Unmarshal(resp, &response) // cannot use resp
	if err != nil {
		return nil, fmt.Errorf("Error while parsing alerts: %v", err)
	}

	alerts := make([]Alert, 0, len(response.Data.Items))
	for _, item := range response.Data.Items {
		rawJSON, err := json.Marshal(item)
		if err != nil {
			continue
		}

		alert := Alert{
			RawJSON: string(rawJSON),
		}

		alerts = append(alerts, alert)
	}
	fmt.Println("A:", alerts, ":A")
	return alerts, nil
}
