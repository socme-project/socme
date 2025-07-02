package routes

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func (r *routerType) RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := r.GetUserFromCookie(c)
		if err != nil {
			c.JSON(401, gin.H{"message": "unauthorized."})
			c.Abort()
			return
		}

		// If RoleMiddleware has no argument: The user need to be connected, all roles allowed
		if len(roles) == 0 {
			c.Next()
			return
		}

		// Else check if the user has one of the roles needed
		role := user.Role
		if slices.Contains(roles, role) {
			c.Next()
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"message": "unauthorized."})
		c.Abort()
	}
}
