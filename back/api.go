package api

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/socme-projects/backend/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/charmbracelet/log"
	"gorm.io/gorm"
)

type API struct {
	Port        string
	Dev         bool
	RefreshRate time.Duration
	Oauth       Oauth
	Logger      *log.Logger
	Router      *gin.Engine
	Db          *gorm.DB
	Domain      string
}

type Oauth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Cfg          *oauth2.Config
}

func NewApi() API {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Prefix:          "SOCme",
	})

	err := godotenv.Load("../.env")
	if err != nil {
		logger.Warn("Error loading .env file in the root folder: ", err)
	}

	if os.Getenv("GITHUB_CLIENT_ID") == "" ||
		os.Getenv("GITHUB_CLIENT_SECRET") == "" ||
		os.Getenv("GITHUB_REDIRECT_URL") == "" {
		logger.Warn("Missing GitHub OAuth environment variables")
	}

	db, err := InitDatabase(utils.GetStringOrDefault(os.Getenv("DB_PATH"), "./socme.db"))
	if err != nil {
		logger.Fatal("Failed to initialize database: ", err)
	}

	refreshRate, err := time.ParseDuration(utils.GetStringOrDefault(os.Getenv("ALERT_RETRIEVAL_INTERVAL"), "5m"))
	if err != nil {
		logger.Fatal("Failed to parse ALERT_RETRIEVAL_INTERVAL: ", err)
	}

	api := API{
		Port:        utils.GetStringOrDefault(os.Getenv("PORT"), "8080"),
		Dev:         utils.GetBoolOrDefault(os.Getenv("DEV"), false),
		RefreshRate: refreshRate,
		Logger:      logger,
		Db:          db,
		Domain:      utils.GetStringOrDefault(os.Getenv("DOMAIN"), "localhost"),
		Oauth: Oauth{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
			Cfg: &oauth2.Config{
				ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
				ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
				RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
				Scopes:       []string{"user:email"},
				Endpoint:     github.Endpoint,
			},
		},
	}

	api.Router = InitRouter(&api)

	return api
}

func (a API) Run() {
	a.Logger.Info("Starting API server on port " + a.Port)
	if err := a.Router.Run(":" + a.Port); err != nil {
		a.Logger.Error("Failed to start API server: ", err)
	}
	a.Logger.Info("API server stopped")
}
