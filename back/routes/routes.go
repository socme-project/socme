package routes

import (
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type routerType struct {
	Db          *gorm.DB
	R           *gin.Engine
	Dev         bool
	Domain      string
	Oauth       Oauth
	RefreshRate time.Duration
	Logger      *log.Logger
}

type Oauth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Cfg          *oauth2.Config
}

func InitRoutes(
	router *gin.Engine,
	db *gorm.DB,
	dev bool,
	logger *log.Logger,
	oauthConfig Oauth,
	domain string,
	refreshRate time.Duration,
) {
	r := routerType{
		R:           router,
		Db:          db,
		Oauth:       oauthConfig,
		Logger:      logger,
		Dev:         dev,
		Domain:      domain,
		RefreshRate: refreshRate,
	}

	go r.UpdateAlerts()

	r.authRoutes()
	r.alertsRoutes()
	r.userRoutes()
	r.clientRoutes()
	r.miscRoutes()
	r.opsmeRoutes()
}
