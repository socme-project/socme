package main

import (
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ClientRoutes defines the routes related to clients
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
			b.Logger.Error("Client not found", "id", id, "error", err)
			c.JSON(http.StatusNotFound, gin.H{"message": "Client not found", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"client": client, "message": "Client retrieved successfully"})
	})

	// POST /client - Create a new client
	b.Router.POST("/client", b.adminMiddleware, func(c *gin.Context) {
		err := model.NewClient(
			b.Db,
			c.PostForm("name"),
			c.PostForm("logo"),
			c.PostForm("wazuhIP"),
			c.PostForm("wazuhPort"),
			c.PostForm("wazuhUsername"),
			c.PostForm("wazuhPassword"),
			c.PostForm("indexerIP"),
			c.PostForm("indexerPort"),
			c.PostForm("indexerUsername"),
			c.PostForm("indexerPassword"),
		)
		if err != nil {
			b.Logger.Error("Error while creating client", "error", err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to create client", "error": err.Error()},
			)
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Client created successfully"})
	})

	// PATCH /client/:id - Update an existing client
	b.Router.PATCH("/client/:id", b.adminMiddleware, func(c *gin.Context) {
		id := c.Param("id")

		var req struct {
			Name            string `json:"name"`
			Logo            string `json:"logo"`
			WazuhIP         string `json:"wazuhIP"`
			WazuhPort       string `json:"wazuhPort"`
			WazuhUsername   string `json:"wazuhUsername"`
			WazuhPassword   string `json:"wazuhPassword"`
			IndexerIP       string `json:"indexerIP"`
			IndexerPort     string `json:"indexerPort"`
			IndexerUsername string `json:"indexerUsername"`
			IndexerPassword string `json:"indexerPassword"`
		}

		// Bind JSON input
		if err := c.ShouldBindJSON(&req); err != nil {
			b.Logger.Error("Invalid request body", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
			return
		}

		// Update client in the database
		err := model.UpdateClient(b.Db, id, req.Name, req.Logo,
			req.WazuhIP, req.WazuhPort, req.WazuhUsername, req.WazuhPassword,
			req.IndexerIP, req.IndexerPort, req.IndexerUsername, req.IndexerPassword)
		if err != nil {
			b.Logger.Error("Failed to update client", "id", id, "error", err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to update client", "error": err.Error()},
			)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Client updated successfully"})
	})
}
