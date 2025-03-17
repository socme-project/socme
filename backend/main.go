package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	wazuhapi "github.com/socme-project/wazuh-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

type Backend struct {
	Port        string
	DbPath      string
	IsProd      bool
	RefreshRate time.Duration
	Wazuh       *wazuhapi.WazuhAPI
	Db          *gorm.DB
	Router      *gin.Engine
	Oauth       Oauth
	Logger      *log.Logger
}

type Oauth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Cfg          *oauth2.Config
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file in the root folder: ", err)
		os.Exit(1)
	}

	interval, _ := time.ParseDuration(os.Getenv("ALERT_RETRIEVAL_INTERVAL"))
	backend := Backend{
		DbPath:      os.Getenv("DB_PATH"),
		Port:        os.Getenv("BACKEND_PORT"),
		IsProd:      os.Getenv("IS_PROD") == "true",
		RefreshRate: interval,
		Logger: log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
			Prefix:          "SOCme",
		}),
		Oauth: Oauth{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		},
		Wazuh: &wazuhapi.WazuhAPI{
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
		},
	}

	backend.Oauth.Cfg = &oauth2.Config{
		ClientID:     backend.Oauth.ClientID,
		ClientSecret: backend.Oauth.ClientSecret,
		RedirectURL:  backend.Oauth.RedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// Default values
	if backend.DbPath == "" {
		backend.DbPath = "data.db"
	}

	if backend.Port == "" {
		backend.Port = "8080"
	}

	if backend.RefreshRate == 0 {
		backend.RefreshRate = 5 * time.Minute
	}

	err = backend.initDB()
	if err != nil {
		backend.Logger.Fatal("Error while initializing the database: ", err)
	}

	if backend.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}
	backend.Router = gin.Default()

	backend.AuthRoutes()
	backend.UserRoutes()
	backend.ClientRoutes()
	backend.AlertRoutes()

	backend.Logger.Info("Starting UpdateAlerts loop")
	go backend.UpdateAlerts()

	// Starting infinite loop to retrieve alerts from Wazuh API
	backend.Logger.Info("Server is launched at http://localhost:" + backend.Port)
	err = backend.Router.Run(":" + backend.Port)
	if err != nil {
		backend.Logger.Fatal("Error while starting the server: " + err.Error())
	}
}
