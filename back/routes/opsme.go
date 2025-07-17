package routes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/socme-project/backend/model"
	"github.com/socme-project/opsme"
)

// isNillArray checks if all items in a slice of errors are nil.
func isNillArray(arr []error) bool {
	for _, item := range arr {
		if item != nil {
			return false
		}
	}
	return true
}

// serializeErrors converts a slice of errors into a slice of strings,
// skipping any nil errors. This is useful for JSON responses.
func serializeErrors(errs []error) []string {
	// Create a slice to hold the error messages.
	// We only include non-nil errors.
	errorStrings := make([]string, 0)
	for _, err := range errs {
		if err != nil {
			errorStrings = append(errorStrings, err.Error())
		}
	}
	return errorStrings
}

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

		errors := UpdateClients(clients)

		fmt.Println("Errors:", errors)
		if !isNillArray(errors) {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Clients were updated. Some could not be updated.",
					"error":   serializeErrors(errors),
				})
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{"message": "Clients updated successfully."},
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

		errors := UpdateClients(
			[]model.Client{*client},
		)

		if !isNillArray(errors) {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to update client.", "error": errors[0].Error()},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{"message": "Client updated successfully."},
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

		fetchOutputs, errors := FetchClients(clients)
		if !isNillArray(errors) {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Some clients could not be fetched.",
					"data":    fetchOutputs,
					"error":   serializeErrors(errors),
				})
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{"message": "Clients fetched successfully.", "data": fetchOutputs},
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
		if !isNillArray(errors) {
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

func UpdateClients(clients []model.Client) []error {
	operator, errs := prepareOpsmeMachines(clients)
	if !isNillArray(errs) {
		return errs
	}

	_, errors := operator.Run(
		"cd /etc/nixos && just pull && nixos-rebuild switch --flake /etc/nixos#node",
	)
	return errors
}

func FetchClients(clients []model.Client) ([]opsme.Output, []error) {
	operator, errs := prepareOpsmeMachines(clients)
	if !isNillArray(errs) {
		return []opsme.Output{}, errs
	}

	results, errors := operator.Run("fastfetch")
	return results, errors
}

func prepareOpsmeMachines(clients []model.Client) (opsme.Operator, []error) {
	homeDir := "/home/socme"
	sshDir := filepath.Join(homeDir, ".ssh")
	knownHostsPath := filepath.Join(sshDir, "known_hosts")

	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		if err := os.MkdirAll(sshDir, 0700); err != nil {
			return opsme.Operator{}, []error{fmt.Errorf("failed to create .ssh directory: %w", err)}
		}
	}

	// Check if the known_hosts file exists. If not, create it.
	if _, err := os.Stat(knownHostsPath); os.IsNotExist(err) {
		file, err := os.Create(knownHostsPath)
		if err != nil {
			return opsme.Operator{}, []error{
				fmt.Errorf("failed to create known_hosts file: %w", err),
			}
		}
		// Close the file immediately and set permissions to 0600 (only owner can read/write).
		file.Close()
		if err := os.Chmod(knownHostsPath, 0600); err != nil {
			return opsme.Operator{}, []error{
				fmt.Errorf("failed to set permissions on known_hosts file: %w", err),
			}
		}
	}

	operator, err := opsme.New(
		true, // this indicates to add to known_hosts file
		5,    // this is the timeout for each operation in seconds
	)
	if err != nil {
		return opsme.Operator{}, []error{err}
	}

	errors := make([]error, len(clients))

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
