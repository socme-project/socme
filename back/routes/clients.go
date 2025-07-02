package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/socme-projects/backend/model"
)

func (r *routerType) clientRoutes() {
	g := r.R.Group("/clients")

	g.GET("", r.RoleMiddleware("admin", "user"), func(c *gin.Context) {
		clients, err := model.GetAllClients(r.Db)
		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to retrieve clients."})
			r.Logger.Error("Failed to retrieve clients", "error", err.Error())
			return
		}
		c.JSON(200, gin.H{"clients": clients})
	})

	g.GET("/:id", r.RoleMiddleware("admin", "user"), func(c *gin.Context) {
		id := c.Param("id")

		client, err := model.GetClientByID(r.Db, id)
		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to retrieve client."})
			r.Logger.Error("Failed to retrieve client", "error", err.Error())
			return
		}
		c.JSON(200, gin.H{"client": client})
	})

	g.DELETE("/:id", r.RoleMiddleware("admin"), func(c *gin.Context) {
		id := c.Param("id")

		err := model.DeleteClient(r.Db, id)
		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to delete client."})
			r.Logger.Error("Failed to delete client", "error", err.Error())
			return
		}
		c.JSON(200, gin.H{"message": "Client deleted successfully."})
		r.Logger.Info("Client deleted successfully", "client_id", id)
	})

	g.POST("", r.RoleMiddleware("admin"), func(c *gin.Context) {
		client, err := model.CreateClient(r.Db,
			c.Query("name"),
			c.Query("logo"),
			c.Query("wazuh_ip"),
			c.Query("wazuh_port"),
			c.Query("wazuh_username"),
			c.Query("wazuh_password"),
			c.Query("indexer_ip"),
			c.Query("indexer_port"),
			c.Query("indexer_username"),
			c.Query("indexer_password"),
		)

		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			r.Logger.Error("Failed to create client", "error", err.Error())
			return
		}

		c.JSON(200, gin.H{"client": client})
		r.Logger.Info("Client created successfully", "client_id", client.ID)
	})

	g.PATCH("/:id", r.RoleMiddleware("admin"), func(c *gin.Context) {
		client, err := model.EditClient(r.Db,
			c.Param("id"),
			c.Query("name"),
			c.Query("logo"),
			c.Query("wazuh_ip"),
			c.Query("wazuh_port"),
			c.Query("wazuh_username"),
			c.Query("wazuh_password"),
			c.Query("indexer_ip"),
			c.Query("indexer_port"),
			c.Query("indexer_username"),
			c.Query("indexer_password"),
		)

		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to edit client."})
			r.Logger.Error("Failed to edit client", "error", err.Error())
			return
		}

		c.JSON(200, gin.H{"client": client})
		r.Logger.Info("Client edited successfully", "client_id", client.ID)
	})

}
