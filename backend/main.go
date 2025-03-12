package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	wazuhapi "github.com/socme-project/wazuh-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

type Backend struct {
	Port                   string
	DbPath                 string
	IsProd                 bool
	AlertRetrievalInterval time.Duration
	Wazuh                  *wazuhapi.WazuhAPI
	Db                     *gorm.DB
	Router                 *gin.Engine
	Oauth                  Oauth
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
		log.Println("Error loading .env file in the root folder: ", err)
	}

	interval, _ := time.ParseDuration(os.Getenv("ALERT_RETRIEVAL_INTERVAL"))
	backend := Backend{
		DbPath:                 os.Getenv("DB_PATH"),
		Port:                   os.Getenv("BACKEND_PORT"),
		IsProd:                 os.Getenv("IS_PROD") == "true",
		AlertRetrievalInterval: interval,
		Oauth: Oauth{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
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

	if backend.AlertRetrievalInterval == 0 {
		backend.AlertRetrievalInterval = 5 * time.Minute
	}

	err = backend.initDB()
	if err != nil {
		log.Fatal("Error while initializing the database: ", err)
	}

	if backend.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}
	backend.Router = gin.Default()

	backend.AuthRoutes()
	backend.UserRoutes()
	backend.ClientRoutes()

	// Starting infinite loop to retrieve alerts from Wazuh API
	log.Println("Server is launched at http://localhost:" + backend.Port)
	err = backend.Router.Run(":" + backend.Port)
	if err != nil {
		log.Fatal("Error while starting the server: ", err)
	}
}
