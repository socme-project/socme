package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/socme-project/backend/model"
	"github.com/socme-project/opsme"
)

func (r *routerType) opsmeRoutes() {
	g := r.R.Group("/opsme")

	g.GET("/update", r.RoleMiddleware("admin"), func(c *gin.Context) {
		clients, err := model.GetAllClients(r.Db)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to fetch clients.", "error": err.Error()},
			)
			return
		}

		updateOutput, errors := UpdateClients(clients)
		if len(errors) > 0 {
			errorMessages := make([]string, len(errors))
			for i, err := range errors {
				errorMessages[i] = err.Error()
			}
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Some clients could not be updated.",
					"error":   errorMessages,
					"data":    updateOutput,
				})
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{"message": "Clients updated successfully.", "data": updateOutput},
		)
	})

	g.GET("/update/:id", r.RoleMiddleware("admin"), func(c *gin.Context) {
		id := c.Param("id")
		client, err := model.GetClientByID(r.Db, id)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to fetch client.", "error": err.Error()},
			)
			return
		}

		updateOutput, errors := UpdateClients(
			[]model.Client{*client},
		)
		if len(errors) > 0 {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to update client.", "error": errors[0].Error()},
			)
			return
		}
		c.JSON(
			http.StatusOK,
			gin.H{"message": "Client updated successfully.", "data": updateOutput},
		)
	})

	g.GET("/fetch", r.RoleMiddleware("admin"), func(c *gin.Context) {
		clients, err := model.GetAllClients(r.Db)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to fetch clients.", "error": err.Error()},
			)
			return
		}

		fetchOutput, errors := FetchClients(clients)
		if len(errors) > 0 {
			errorMessages := make([]string, len(errors))
			for i, err := range errors {
				errorMessages[i] = err.Error()
			}

			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Some clients could not be fetched.",
					"error":   errorMessages,
					"data":    fetchOutput,
				})
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{"message": "Clients fetched successfully.", "data": fetchOutput},
		)
	})

	g.GET("/fetch/:id", r.RoleMiddleware("admin"), func(c *gin.Context) {
		id := c.Param("id")
		client, err := model.GetClientByID(r.Db, id)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to fetch client.", "error": err.Error()},
			)
			return
		}

		fetchOutput, errors := FetchClients(
			[]model.Client{*client},
		)
		if len(errors) > 0 {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to fetch client.", "error": errors[0].Error()},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{"message": "Client fetched successfully.", "data": fetchOutput},
		)
	})
}

func UpdateClients(clients []model.Client) ([]opsme.Output, []error) {
	operator, err := prepareOpsmeMachines(clients)
	if err != nil {
		return []opsme.Output{}, err
	}

	results, errors := operator.Run(
		"cd /etc/nixos && just pull && nixos-rebuild switch --flake /etc/nixos#node",
	)
	return results, errors
}

func FetchClients(clients []model.Client) ([]opsme.Output, []error) {
	operator, err := prepareOpsmeMachines(clients)
	if err != nil {
		return []opsme.Output{}, err
	}

	results, errors := operator.Run("fastfetch")
	return results, errors
}

func prepareOpsmeMachines(clients []model.Client) (opsme.Operator, []error) {
	operator, err := opsme.New(
		true, // this indicates to add to known_hosts file
		5,    // this is the timeout for each operation in seconds
	)
	if err != nil {
		return opsme.Operator{}, []error{err}
	}

	errors := make([]error, 0, len(clients))

	for i, client := range clients {
		SshPortInt, err := strconv.Atoi(client.SshPort)
		if err != nil {
			errors[i] = err
			continue
		}

		machine, err := operator.NewMachine(
			client.Name,
			client.SshUsername,
			client.Host,
			SshPortInt,
		)
		if err != nil {
			errors[i] = err
			continue
		}

		err = machine.WithPasswordAuth(client.SshPassword)
		if err != nil {
			errors[i] = err
			continue
		}
	}

	return operator, errors
}
