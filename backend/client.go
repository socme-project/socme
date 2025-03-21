package main

import (
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (b *Backend) ClientRoutes() {
	// GET /clients - List all clients
	b.Router.GET("/clients", b.userMiddleware, func(c *gin.Context) {
		clients := model.GetAllClients(b.Db)
		c.JSON(http.StatusOK, gin.H{"clients": clients})
	})

	// GET /client/:id - Get a client by ID
	b.Router.GET("/client/:id", b.userMiddleware, func(c *gin.Context) {
		id := c.Param("id")
		client, err := model.GetClientByID(b.Db, id)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"message": "Client not found.", "error": err.Error()},
			)
			return
		}
		c.JSON(http.StatusOK, gin.H{"client": client, "message": "Client retrieved."})
	})

	// POST /client - Create a new client
	b.Router.POST("/client", b.adminMiddleware, func(c *gin.Context) {
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
			c.JSON(
				http.StatusBadRequest,
				gin.H{"message": "Error while creating the client.", "error": err.Error()},
			)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Client created."})
	})
}
