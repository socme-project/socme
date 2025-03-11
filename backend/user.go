package main

import (
	"backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
	certfr_scrapping "github.com/socme-project/cert-fr-scrapping"
)

func (b *Backend) UserRoutes() {
	b.Router.GET("/refreshUser", b.userMiddleware, func(c *gin.Context) {
		user, err := model.GetUserByToken(b.Db, c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Error getting user with token"})
		}
		c.JSON(http.StatusOK, gin.H{"message": "You are connected", "user": user})
	})

	b.Router.GET("/users/list", b.userMiddleware, func(c *gin.Context) {
		users := []model.User{}
		b.Db.Find(&users)
		c.JSON(http.StatusOK, gin.H{"users": users})
	})

	b.Router.GET("/certfr", b.userMiddleware, func(c *gin.Context) {
		alerts := certfr_scrapping.CollectAlert()
		cti := certfr_scrapping.CollectCti()
		avis := certfr_scrapping.CollectAvis()
		c.JSON(http.StatusOK, gin.H{"alerts": alerts, "cti": cti, "avis": avis})
	})

	b.Router.GET("/user/change-role", b.adminMiddleware, func(c *gin.Context) {
		id := c.Query("id")
		newRole := c.Query("role")
		err := model.UpdateUserRole(b.Db, id, newRole)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
	})

	b.Router.GET("/user/revoke-session", b.adminMiddleware, func(c *gin.Context) {
		id := c.Query("id")
		var session model.Session
		if err := b.Db.First(&session, "ID = ?", id).Error; err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{"message": "Session not found for this user-id"},
			)
		}
		b.Db.Delete(&session)
		c.JSON(http.StatusOK, gin.H{"message": "Session revoked"})
	})

	b.Router.GET("/user/delete", b.adminMiddleware, func(c *gin.Context) {
		id := c.Query("id")
		user, err := model.GetUserByID(b.Db, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		b.Db.Delete(&user)
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	})
}
