package headscale

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	hsClient "github.com/hibare/headscale-client-go"
)

type HeadscaleServer struct {
	ServerUrl string
	ApiToken  string
}

func (HeadscaleServer) NewClient() {

}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug) // Optional

	serverUrl := os.Getenv("HS_SERVER_URL")
	apiToken := os.Getenv("HS_API_TOKEN")

	client, err := hsClient.NewClient(serverUrl, apiToken, hsClient.HeadscaleClientOptions{})
	if err != nil {
		panic(err)
	}

	nodes, err := client.Nodes().List(context.Background())
	if err != nil {
		panic(err)
	}

	for _, node := range nodes.Nodes {
		fmt.Printf("Node: %s\n", node.Name)
	}
}
