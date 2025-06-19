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
	"gorm.io/gorm/logger"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().
			Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

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
	isProd := os.Getenv("IS_PROD") == "true"

	if !isProd {
		err := godotenv.Load("../.env")
		if err != nil && !os.IsNotExist(err) {
			fmt.Println("Error loading .env file in the root folder: ", err)
			os.Exit(1)
		}
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data.db"
	}

	interval, err := time.ParseDuration(os.Getenv("ALERT_RETRIEVAL_INTERVAL"))
	if err != nil {
		interval = 5 * time.Minute
	}

	backend := Backend{
		Port:        os.Getenv("BACKEND_PORT"),
		IsProd:      isProd,
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

	if backend.Port == "" {
		backend.Port = "8080"
	}

	if backend.RefreshRate == 0 {
		backend.RefreshRate = 5 * time.Minute
	}

	err = backend.initDB(dbPath)
	if err != nil {
		backend.Logger.Fatal("Error while initializing the database: ", err)
	}

	if backend.IsProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	backend.Router = gin.Default()

	backend.Router.Use(CORSMiddleware())

	backend.AuthRoutes()
	backend.UserRoutes()
	backend.ClientRoutes()
	backend.AlertRoutes()
	backend.MiscRoutes()

	backend.Logger.Info("Starting UpdateAlerts loop")
	go backend.UpdateAlerts()

	backend.Logger.Info("Server is launched at http://localhost:" + backend.Port)
	err = backend.Router.Run(":" + backend.Port)
	if err != nil {
		backend.Logger.Fatal("Error while starting the server: " + err.Error())
	}
}

func (backend *Backend) initDB(dbPath string) error {
	var err error
	// TODO: Add custom logger
	backend.Db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
