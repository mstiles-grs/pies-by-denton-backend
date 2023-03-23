package middleware

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"time"

)

func UpdateSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Execute other middlewares and handlers first
		c.Next()

		// Get the user ID from the request context
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Generate a new session token
		newSessionToken := generateSessionToken()

		// Set the new session token as a cookie
		c.SetCookie("session_token", newSessionToken, 3600, "", "", false, true)

		// Connect to Redis
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		ctx := context.Background()

		// Delete the old session token from Redis
		oldSessionToken, _ := c.Cookie("session_token")
		rdb.Del(ctx, oldSessionToken)

		// Set the new session token in Redis, associated with the user ID
		rdb.Set(ctx, newSessionToken, userID, time.Hour*24)
	}
}

func generateSessionToken() string {
	sessionToken := uuid.New().String()

	return sessionToken
}


