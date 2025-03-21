package main

import (
	"backend/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (b *Backend) UserRoutes() {
	// GET /users - List all users
	b.Router.GET("/users", b.adminMiddleware, func(c *gin.Context) {
		users := []model.User{}
		b.Db.Find(&users)
		c.JSON(http.StatusOK, gin.H{"users": users, "message": "List of users returned"})
	})

	// PATCH /users/:id/role - Update a user's role
	b.Router.PATCH("/users/:id/role", b.adminMiddleware, func(c *gin.Context) {
		id := c.Param("id")
		newRole := c.Query("role")
		err := model.UpdateUserRole(b.Db, id, newRole)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
	})

	// DELETE /users/:id/session - Revoke a user's session
	b.Router.DELETE("/users/:id/session", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user-id"})
			return
		}
		user, err := model.GetUserFromRequest(b.Db, c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		if user.Role != "admin" || user.ID != uint(id) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		var session model.Session
		if err := b.Db.First(&session, "ID = ?", id).Error; err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"message": "Session not found for this user-id"},
			)
			return
		}
		b.Db.Delete(&session)
		c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
	})

	// DELETE /users/:id - Delete a user
	b.Router.DELETE("/users/:id", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user-id"})
			return
		}
		user, err := model.GetUserFromRequest(b.Db, c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		if user.Role != "admin" || user.ID != uint(id) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		user, err = model.GetUserByID(b.Db, strconv.FormatUint(id, 10))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		b.Db.Delete(&user)
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	})
}
