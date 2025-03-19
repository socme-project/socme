package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	certfr_scrapping "github.com/socme-project/cert-fr-scrapping"
)

// MiscRoutes defines the routes for the miscellaneous functions, including the certfr route (for alerts, cti, and more)
func (b *Backend) MiscRoutes() {
	b.Router.GET("/certfr", b.userMiddleware, func(c *gin.Context) {
		alerts := certfr_scrapping.CollectAlert()
		cti := certfr_scrapping.CollectCti()
		avis := certfr_scrapping.CollectAvis()
		c.JSON(
			http.StatusOK,
			gin.H{"alerts": alerts, "cti": cti, "avis": avis, "message": "List of alerts returned"},
		)
	})
}
