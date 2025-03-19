package main

import (
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ClientRoutes defines the routes for the client, including the client list, client by ID, and new client
func (b *Backend) ClientRoutes() {
	b.Router.GET("/clients/list", b.userMiddleware, func(c *gin.Context) {
		clients := model.GetAllClients(b.Db)
		c.JSON(http.StatusOK, gin.H{"clients": clients})
	})

	b.Router.GET("/client/:id", b.userMiddleware, func(c *gin.Context) {
		id := c.Param("id")
		client, err := model.GetClientByID(b.Db, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"client": client, "message": "Client retrieved"})
	})

	b.Router.GET("/client/new", b.userMiddleware, func(c *gin.Context) {
		err := model.NewClient(
			b.Db,
			c.Query("name"),
			c.Query("logo"),
			c.Query("wazuhIP"),
			c.Query("wazuhPort"),
			c.Query("wazuhUsername"),
			c.Query("wazuhPassword"),
			c.Query("indexerIP"),
			c.Query("indexerPort"),
			c.Query("indexerUsername"),
			c.Query("indexerPassword"),
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "Client created"})
	})
}
