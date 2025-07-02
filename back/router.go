package api

import (
	"github.com/gin-gonic/gin"
	"github.com/socme-projects/backend/routes"
)

func InitRouter(api *API) *gin.Engine {
	router := gin.Default()

	router.Use(CORSMiddleware())

	// Routes
	routes.InitRoutes(
		router,
		api.Db,
		api.Dev,
		api.Logger,
		routes.Oauth{
			ClientID:     api.Oauth.ClientID,
			ClientSecret: api.Oauth.ClientSecret,
			RedirectURL:  api.Oauth.RedirectURL,
			Cfg:          api.Oauth.Cfg,
		},
		api.Token,
		api.Domain,
	)

	return router
}
