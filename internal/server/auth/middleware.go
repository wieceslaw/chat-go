package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const userKey = "user"

type AuthMiddleware struct {
	Service UserService
}

func NewAuthMiddleware(service UserService) *AuthMiddleware {
	return &AuthMiddleware{
		Service: service,
	}
}

func GetUser(c *gin.Context) (*User, error) {
	user, ok := c.Get(userKey)
	println(user)
	if ok {
		return user.(*User), nil
	}
	return nil, errors.New("No user in context")
}

func (am *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			return
		}

		user, err := am.Service.ValidateToken(c.Request.Context(), AuthToken(token))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			return
		}

		c.Set(userKey, user)
		c.Next()
	}
}
