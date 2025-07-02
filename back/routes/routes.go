package routes

import (
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type routerType struct {
	Db     *gorm.DB
	R      *gin.Engine
	Dev    bool
	Token  []byte
	Domain string
	Oauth  Oauth
	Logger *log.Logger
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
	token []byte,
	domain string,
) {
	r := routerType{
		R:      router,
		Db:     db,
		Oauth:  oauthConfig,
		Logger: logger,
		Token:  token,
		Dev:    dev,
		Domain: domain,
	}

	r.authRoutes()
	r.userRoutes()
	r.clientRoutes()
}
