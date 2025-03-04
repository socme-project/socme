package main

import (
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (b *Backend) ClientRoutes() {
	b.Router.GET("/clients/list", b.userMiddleware, func(c *gin.Context) {
		clients := model.GetAllClients(b.Db)
		c.JSON(http.StatusOK, gin.H{"clients": clients})
	})

	b.Router.GET("/client/new", b.userMiddleware, func(c *gin.Context) {
		err := model.NewClient(
			b.Db,
			c.Query("name"),
			c.Query("logo"),
			c.Query("artemisIP"),
			c.Query("artemisPassword"),
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "Client created"})
	})
}
