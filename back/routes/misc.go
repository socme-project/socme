package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	certfr_scraping "github.com/socme-project/cert-fr-scraping"
)

func (r *routerType) miscRoutes() {
	g := r.R.Group("/misc")

	g.GET("certfr", func(c *gin.Context) {
		alerts := certfr_scraping.CollectAlert()
		cti := certfr_scraping.CollectCti()
		avis := certfr_scraping.CollectAvis()
		c.JSON(
			http.StatusOK,
			gin.H{
				"alerts":  alerts,
				"cti":     cti,
				"avis":    avis,
				"message": "List of alerts returned.",
			},
		)
	})
}
