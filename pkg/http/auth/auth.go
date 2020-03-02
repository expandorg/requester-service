package auth

import (
	"errors"

	"github.com/expandorg/requester-service/pkg/backend"
	"github.com/gemsorg/svc-kit/http/api"
	"github.com/gin-gonic/gin"
)

const (
	jwtKey    = "JWT"
	userIDKey = "user-id"
)

// Middleware for auth
func Middleware(c *gin.Context) error {
	userID, err := GetRequestUserID(c)
	if err != nil {
		return err
	}
	if userID == 0 {
		return api.Forbidden(UnableToAuthErr{})
	}
	c.Set(userIDKey, userID)
	return nil
}

// Get Request User ID
func GetRequestUserID(c *gin.Context) (uint64, error) {
	cookie, err := GetJWTCookie(c)
	if err != nil {
		return 0, api.Forbidden(err)
	}

	userID, err := backend.GetAuth(cookie)
	if err != nil {
		if _, ok := err.(UnableToAuthErr); ok {
			return 0, api.Forbidden(err)
		}
		return 0, err
	}
	return userID, nil
}

// GetJWTCookie from gin request
func GetJWTCookie(c *gin.Context) (string, error) {
	cookie, err := c.Cookie(jwtKey)
	if err != nil {
		return "", err
	}
	return cookie, nil
}

// GetUserID from gin request
func GetUserID(c *gin.Context) (uint64, error) {
	userID, ok := c.Get(userIDKey)
	if !ok {
		return 0, errors.New("Could not identify user")
	}
	return userID.(uint64), nil
}

type UnableToAuthErr struct{}

func (err UnableToAuthErr) Error() string {
	return "Unable to authenticate user"
}
