package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/socme-project/backend/model"
)

func (r *routerType) userRoutes() {
	g := r.R.Group("/users")

	// List all users, only for admin
	g.GET("", r.RoleMiddleware("admin"), func(c *gin.Context) {
		users, err := model.GetAllUsers(r.Db)
		if err != nil {
			c.JSON(500, gin.H{"message": "Failed to retrieve users."})
			r.Logger.Error("Failed to retrieve users", "error", err.Error())
			return
		}
		c.JSON(200, gin.H{"users": users})
	})

	g.PATCH("/:id/role", r.RoleMiddleware("admin"), func(c *gin.Context) {
		id := c.Param("id")
		newRole := c.Query("role")
		err := model.EditUserRole(r.Db, id, newRole)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"message": "Error while updating the role.", "error": err.Error()},
			)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Role updated."})
	})

	g.DELETE("/:id/session", r.RoleMiddleware("admin", "user"), func(c *gin.Context) {
		user, err := r.GetUserFromCookie(c)
		if err != nil {
			c.JSON(401, gin.H{"message": "unauthorized."})
			c.Abort()
			return
		}
		if user.Role != "admin" && user.ID != c.Param("id") {
			c.JSON(401, gin.H{"message": "unauthorized."})
			c.Abort()
			return
		}
		token, err := GetTokenFromCookie(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to log out."})
			r.Logger.Error("Failed to log out", "error", err.Error())
			return
		}
		model.DeleteSession(r.Db, token)
	})

	g.DELETE("/:id", r.RoleMiddleware("admin", "user"), func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User ID is required."})
			return
		}
		user, err := r.GetUserFromCookie(c)
		if err != nil {
			c.JSON(401, gin.H{"message": "unauthorized."})
			c.Abort()
			return
		}
		if user.Role != "admin" && user.ID != id {
			c.JSON(http.StatusForbidden, gin.H{"message": "unauthorized."})
			return
		}

		_ = model.DeleteSession(r.Db, user.ID)

		err = model.DeleteUser(r.Db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user."})
			r.Logger.Error("Failed to delete user", "error", err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted."})
	})
}
