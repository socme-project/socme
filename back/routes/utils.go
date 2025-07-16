package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/socme-project/backend/model"
)

func GetTokenFromCookie(c *gin.Context) (string, error) {
	token, err := c.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		return "", errors.New("token not found")
	}
	if token == "" {
		return "", errors.New("token empty")
	}
	return token, nil
}

func (r *routerType) GetUserFromCookie(c *gin.Context) (model.User, error) {
	token, err := GetTokenFromCookie(c)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get token from cookie: %w", err)
	}
	u, err := r.GetUserFromToken(token)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user from token: %w", err)
	}

	return *u, nil
}

func (r *routerType) GetUserFromToken(tokenString string) (*model.User, error) {
	s, err := model.GetSession(r.Db, tokenString)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	return &s.User, nil
}

func parseGithubUser(body io.ReadCloser) (githubID string, githubLogin string, err error) {
	var githubUser map[string]any
	err = json.NewDecoder(body).Decode(&githubUser)
	if err != nil {
		err = errors.New("failed to decode the answer from Github")
		return
	}

	id, ok := githubUser["id"].(float64)
	if !ok {
		idStr, isStr := githubUser["id"].(string)
		if isStr {
			githubID = idStr
		} else {
			err = errors.New("failed to parse GitHub user ID: unexpected type")
			return
		}
	} else {
		githubID = strconv.FormatFloat(id, 'f', 0, 64)
	}

	githubLogin, ok = githubUser["login"].(string)
	if !ok {
		err = errors.New("failed to parse GitHub user login")
		return
	}

	return
}
