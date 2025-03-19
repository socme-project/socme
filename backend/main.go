package main

import (
	"backend/model"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Backend struct {
	Port        string
	IsProd      bool
	RefreshRate time.Duration
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

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data.db"
	}

	interval, _ := time.ParseDuration(os.Getenv("ALERT_RETRIEVAL_INTERVAL"))
	backend := Backend{
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
	}

	backend.Oauth.Cfg = &oauth2.Config{
		ClientID:     backend.Oauth.ClientID,
		ClientSecret: backend.Oauth.ClientSecret,
		RedirectURL:  backend.Oauth.RedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// Default values
	if backend.Port == "" {
		backend.Port = "8080"
	}

	if backend.RefreshRate == 0 {
		backend.RefreshRate = 5 * time.Minute
	}

	// Start
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
	backend.MiscRoutes()

	backend.Logger.Info("Starting UpdateAlerts loop")
	go backend.UpdateAlerts()

	// Starting infinite loop to retrieve alerts from Wazuh API
	backend.Logger.Info("Server is launched at http://localhost:" + backend.Port)
	err = backend.Router.Run(":" + backend.Port)
	if err != nil {
		backend.Logger.Fatal("Error while starting the server: " + err.Error())
	}
}

func (backend *Backend) initDB() error {
	var err error
	backend.Db, err = gorm.Open(sqlite.Open(backend.DbPath), &gorm.Config{
		// TODO: Custom logger
		Logger: nil,
	})
	if err != nil {
		return err
	}

	err = backend.Db.AutoMigrate(&model.User{}, &model.Session{}, &model.Client{}, &model.Alert{})
	if err != nil {
		return err
	}

	return nil
}
