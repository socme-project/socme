package main

import (
	"backend/model"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (b *Backend) AuthRoutes() {
	b.Router.GET("/auth/github", b.authFunc)
	b.Router.GET("/auth/callback", b.authCallbackFunc)
	b.Router.GET("/auth/refresh", func(c *gin.Context) {
		user, err := model.GetUserByToken(b.Db, c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		}
		c.JSON(http.StatusOK, gin.H{"message": "User returned", "user": user})
	})
}

func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "default_state"
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (b *Backend) authFunc(c *gin.Context) {
	state := generateState()
	if b.IsProd {
		c.SetCookie("oauth_state", state, 300, "/", "", true, true)
	} else {
		c.SetCookie("oauth_state", state, 300, "/", "", false, true)
	}

	url := b.Oauth.Cfg.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

func (b *Backend) authCallbackFunc(c *gin.Context) {
	storedState, err := c.Cookie("oauth_state")

	if err != nil || c.Query("state") != storedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state"})
		return
	}

	code := c.Query("code")
	token, err := b.Oauth.Cfg.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to exchange token"})
		return
	}

	client := b.Oauth.Cfg.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Failed to get user from Github"},
		)
		return
	}
	defer resp.Body.Close()

	var githubUser map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&githubUser)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Failed to decode the answer from Github"},
		)
		return
	}

	// Make the first user admin
	var count int64
	b.Db.Model(&model.User{}).Count(&count)

	var user model.User
	b.Db.FirstOrCreate(&user, model.User{
		GitHubID: json.Number(fmt.Sprintf("%.0f", githubUser["id"].(float64))).String(),
		Username: githubUser["login"].(string),
		Avatar:   githubUser["avatar_url"].(string),
	})
	if user.Role == "" {
		if count == 0 {
			user.Role = "admin"
		} else {
			user.Role = "guest" // Default role
		}
	}
	b.Db.Save(&user)

	sessionToken := token.AccessToken
	session := model.Session{
		Token: sessionToken, ID: user.ID,
		ExpiresAt: time.Now().Add(10 * time.Hour),
	}

	// Update or create the session
	if b.Db.Model(&session).Where("id = ?", user.ID).Updates(&user).RowsAffected == 0 {
		b.Db.Create(&session)
	}
	b.Db.Save(&session)

	c.JSON(http.StatusOK, gin.H{"token": sessionToken, "user": user})
}
