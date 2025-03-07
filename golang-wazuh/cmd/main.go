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

	fmt.Println(apiVer)
	alerts, err := wazuh.GetAlerts()
	if err != nil {
		panic(err)
	}
	fmt.Println(alerts)
}
