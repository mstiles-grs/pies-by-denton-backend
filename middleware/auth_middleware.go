package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// TODO: Use the session ID to get user information and pass it downstream.

		c.Next()
	}
}