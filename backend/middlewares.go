package main

import (
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (b *Backend) authMiddleware(c *gin.Context) {
	_, err := model.GetUserFromRequest(b.Db, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func (b *Backend) userMiddleware(c *gin.Context) {
	user, err := model.GetUserFromRequest(b.Db, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	if user.Role != "user" && user.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}

func (b *Backend) adminMiddleware(c *gin.Context) {
	user, err := model.GetUserFromRequest(b.Db, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	if user.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
	c.Next()
}
