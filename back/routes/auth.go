package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/socme-projects/backend/model"
	"github.com/socme-projects/backend/utils"
)

var SESSION_COOKIE_NAME = "session"

func (r *routerType) authRoutes() {
	g := r.R.Group("/auth")

	g.GET("/refresh", r.RoleMiddleware(), func(c *gin.Context) {
		user, err := r.GetUserFromCookie(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized."})
			r.Logger.Error("Failed to refresh session", "error", err.Error())
			return
		}
		session, err := model.CreateSession(r.Db, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to refresh session."})
			r.Logger.Error("Failed to create session", "error", err.Error())
			return
		}
		c.SetCookie(SESSION_COOKIE_NAME, session.Token, 45000, "/", r.Domain, !r.Dev, false)
		c.JSON(http.StatusOK, gin.H{"user": user})
		r.Logger.Info("Session refreshed successfully", "user_id", user.ID)
	})

	g.GET("/logout", r.RoleMiddleware(), func(c *gin.Context) {
		token, err := c.Cookie(SESSION_COOKIE_NAME)
		if err != nil || token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to log out."})
			r.Logger.Error("Failed to log out", "error", err.Error(), "token", token)
			return
		}
		err = model.DeleteSession(r.Db, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to log out."})
			r.Logger.Error("Failed to log out", "error", err.Error())
			return
		}
		c.SetCookie(SESSION_COOKIE_NAME, "", -1, "/", r.Domain, !r.Dev, false)
		c.SetCookie("oauth_state", "", -1, "/", r.Domain, !r.Dev, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully."})
		r.Logger.Info("User logged out successfully")
	})

	g.GET("/github/login", func(c *gin.Context) {
		state := utils.GenerateState()
		c.SetCookie("oauth_state", state, 300, "/", r.Domain, !r.Dev, true)

		url := r.Oauth.Cfg.AuthCodeURL(state)
		c.Redirect(http.StatusFound, url)
	})

	g.GET("/github/callback", func(c *gin.Context) {
		storedState, err := c.Cookie("oauth_state")

		if err != nil || c.Query("state") != storedState {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid state."})
			r.Logger.Error("Invalid state in OAuth callback", "error", err)
			return
		}

		code := c.Query("code") // Get the code from github
		githubToken, err := r.Oauth.Cfg.Exchange(c.Request.Context(), code)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to exchange token."},
			)
			r.Logger.Error("Failed to exchange token", "error", err.Error())
			return
		}

		client := r.Oauth.Cfg.Client(c.Request.Context(), githubToken)
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"message": "Failed to get user from Github."},
			)
			r.Logger.Error("Failed to get user from Github", "error", err.Error())
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		githubID, githubLogin, err := parseGithubUser(resp.Body)

		//
		// Get the user
		//

		u, err := model.GetUserByGithubID(r.Db, githubID)
		if err != nil {
			u, err = model.CreateUser(r.Db, githubLogin, githubID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user."})
				r.Logger.Error("Failed to create user", "error", err.Error())
				return
			}
		}

		sessionToken, err := model.CreateSession(r.Db, *u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create session."})
			r.Logger.Error("Failed to create session", "error", err.Error())
			return
		}
		timestamp := sessionToken.Exp.Unix()
		timestampStr := strconv.FormatInt(timestamp, 10)
		c.SetCookie(SESSION_COOKIE_NAME, sessionToken.Token, 45000, "/", r.Domain, !r.Dev, false)
		c.SetCookie("exp", timestampStr, 45000, "/", r.Domain, !r.Dev, false)
		c.JSON(http.StatusOK, gin.H{"user": *u})
	})
}
