package main

import (
	"fmt"

	wazuhapi "github.com/opensoc-paris/wazuh"
)

func main() {
	wazuh := wazuhapi.WazuhAPI{
		Host:     "10.8.178.20",
		Port:     "55000",
		Username: "admin",
		Password: "HMthisismys3cr3tP5ssword34a;",
		Indexer: wazuhapi.Indexer{
			Username: "admin",
			Password: "HMthisismys3cr3tP5ssword34a;",
			Host:     "10.8.178.20",
			Port:     "9200",
		},
		Insecure: true,
	}

	err := wazuh.RefreshToken()
	if err != nil {
		panic(err)
	}

	apiVer, err := wazuh.GetApiVersion()
	if err != nil {
		panic(err)
	}
	lastAlertId := "0"
	fmt.Println(apiVer)
	alerts, err := wazuh.GetAlerts(lastAlertId)
	if err != nil {
		panic(err)
	}

	// debug
	fmt.Println(len(alerts))
}
