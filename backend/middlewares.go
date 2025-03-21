package main

import (
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// authMiddleware checks if the user is authenticated (guest or user or admin)
// func (b *Backend) authMiddleware(c *gin.Context) {
// 	_, err := model.GetUserFromRequest(b.Db, c)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
// 		c.Abort()
// 		return
// 	}
// 	c.Next()
// }

// userMiddleware checks if the user is authenticated as a user or admin
func (b *Backend) userMiddleware(c *gin.Context) {
	user, err := model.GetUserFromRequest(b.Db, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized."})
		c.Abort()
		return
	}
	if user.Role != "user" && user.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized."})
		c.Abort()
		return
	}
	c.Next()
}

// adminMiddleware checks if the user is authenticated as an admin
func (b *Backend) adminMiddleware(c *gin.Context) {
	user, err := model.GetUserFromRequest(b.Db, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized."})
		c.Abort()
		return
	}
	if user.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized."})
		c.Abort()
		return
	}
	c.Next()
}
