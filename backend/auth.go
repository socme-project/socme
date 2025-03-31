package main

import (
	"backend/model"
	"backend/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// AuthRoutes sets up routes for both GitHub and Google OAuth authentication.
func (b *Backend) AuthRoutes() {
	auth := b.Router.Group("/auth")
	auth.GET("/github/login", b.GithubLogin)
	auth.GET("/github/callback", b.GithubCallback)
	auth.GET("/google/login", b.GoogleLogin)
	auth.GET("/google/callback", b.GoogleCallback)

	auth.GET("/refresh", func(c *gin.Context) {
		user, err := model.GetUserByToken(b.Db, c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized.", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User returned.", "user": user})
	})
}

// GithubLogin initiates the GitHub OAuth flow.
func (b *Backend) GithubLogin(c *gin.Context) {
	state := utils.GenerateState()
	if b.IsProd {
		c.SetCookie("oauth_state_github", state, 300, "/", "", true, true)
	} else {
		c.SetCookie("oauth_state_github", state, 300, "/", "", false, true)
	}
	url := b.Oauth.GitHubCfg.AuthCodeURL(state)
	c.Redirect(http.StatusFound, url)
}

// GoogleLogin initiates the Google OAuth flow.
func (b *Backend) GoogleLogin(c *gin.Context) {
	state := utils.GenerateState()
	if b.IsProd {
		c.SetCookie("oauth_state_google", state, 300, "/", "", true, true)
	} else {
		c.SetCookie("oauth_state_google", state, 300, "/", "", false, true)
	}
	url := b.Oauth.GoogleCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

// GithubCallback handles GitHub OAuth callback.
func (b *Backend) GithubCallback(c *gin.Context) {
	storedState, err := c.Cookie("oauth_state_github")
	if err != nil || c.Query("state") != storedState {
		errMsg := "Invalid state."
		if err != nil {
			errMsg = fmt.Sprintf("Invalid state: %s", err.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
		return
	}

	code := c.Query("code")
	token, err := b.Oauth.GitHubCfg.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Failed to exchange token.", "error": err.Error()},
		)
		return
	}

	client := b.Oauth.GitHubCfg.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Failed to get user from GitHub.", "error": err.Error()},
		)
		return
	}
	defer resp.Body.Close()

	var githubUser map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Failed to decode GitHub response.", "error": err.Error()},
		)
		return
	}

	githubID := json.Number(fmt.Sprintf("%.0f", githubUser["id"].(float64))).String()
	username := githubUser["login"].(string)

	// Check if a user with this username exists but doesn't have a GitHub ID.
	var existingUser model.User
	if err := b.Db.Where("username = ? AND github_id IS NULL", username).First(&existingUser).Error; err == nil {
		username += " (github)"
	}

	var user model.User
	b.Db.FirstOrCreate(&user, model.User{
		GitHubID: githubID,
		Username: username,
		Avatar:   githubUser["avatar_url"].(string),
	})

	// Set role if not already assigned.
	var count int64
	b.Db.Model(&model.User{}).Count(&count)
	if user.Role == "" {
		if count == 0 {
			user.Role = "admin"
		} else {
			user.Role = "guest"
		}
	}
	b.Db.Save(&user)

	sessionToken := token.AccessToken
	session := model.Session{
		Token:     sessionToken,
		ID:        user.ID,
		ExpiresAt: time.Now().Add(10 * time.Hour),
	}
	if b.Db.Model(&session).Where("id = ?", user.ID).Updates(&user).RowsAffected == 0 {
		b.Db.Create(&session)
	}
	b.Db.Save(&session)

	c.JSON(http.StatusOK, gin.H{"token": sessionToken, "user": user})
}

// GoogleCallback handles Google OAuth callback.
func (b *Backend) GoogleCallback(c *gin.Context) {
	storedState, err := c.Cookie("oauth_state_google")
	if err != nil || c.Query("state") != storedState {
		errMsg := "Invalid state."
		if err != nil {
			errMsg = fmt.Sprintf("Invalid state: %s", err.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": errMsg})
		return
	}

	code := c.Query("code")
	token, err := b.Oauth.GoogleCfg.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Failed to exchange token.", "error": err.Error()},
		)
		return
	}

	client := b.Oauth.GoogleCfg.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Failed to get user from Google.", "error": err.Error()},
		)
		return
	}
	defer resp.Body.Close()

	var googleUser map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"message": "Failed to decode Google response.", "error": err.Error()},
		)
		return
	}

	// Extract necessary information from the Google response.
	googleID, ok := googleUser["sub"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid Google user ID"})
		return
	}
	username, _ := googleUser["name"].(string)
	avatar, _ := googleUser["picture"].(string)

	// Check if a user with this username exists but doesn't have a Google ID.
	var existingUser model.User
	if err := b.Db.Where("username = ? AND google_id IS NULL", username).First(&existingUser).Error; err == nil {
		username += " (google)"
	}

	var user model.User
	b.Db.FirstOrCreate(&user, model.User{
		GoogleID: googleID,
		Username: username,
		Avatar:   avatar,
	})

	// Set role if not already assigned.
	var count int64
	b.Db.Model(&model.User{}).Count(&count)
	if user.Role == "" {
		if count == 0 {
			user.Role = "admin"
		} else {
			user.Role = "guest"
		}
	}
	b.Db.Save(&user)

	sessionToken := token.AccessToken
	session := model.Session{
		Token:     sessionToken,
		ID:        user.ID,
		ExpiresAt: time.Now().Add(10 * time.Hour),
	}
	if b.Db.Model(&session).Where("id = ?", user.ID).Updates(&user).RowsAffected == 0 {
		b.Db.Create(&session)
	}
	b.Db.Save(&session)

	c.JSON(http.StatusOK, gin.H{"token": sessionToken, "user": user})
}
