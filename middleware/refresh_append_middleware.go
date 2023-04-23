package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/mstiles-grs/pies-by-denton-backend/database"
)



type JsonResponse struct {
    Data         json.RawMessage `json:"data"`
    SessionToken string          `json:"session_token,omitempty"`
}

type CustomResponseWriter struct {
    gin.ResponseWriter
    body          []byte
    newSessionToken string
    jsonResponse  JsonResponse
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
    w.body = b
    err := json.Unmarshal(b, &w.jsonResponse.Data)
    if err != nil {
        return 0, err
    }
    return len(b), nil
}

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Find the user ID associated with the session token
		sessionToken := c.GetHeader("Authorization")
		fmt.Printf("sessionToken: sessionToken = %s", sessionToken)

		if len(sessionToken) > 7 && strings.HasPrefix(sessionToken, "Bearer ") {
			sessionToken = sessionToken[7:]
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid session token"})
			return
		}

		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		userID, err := rdb.Get(database.DbContext, sessionToken).Result()
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid session token"})
			return
		}

		// Remove the old session token from Redis
		rdb.Unlink(database.DbContext, sessionToken)

		// Generate a new session token and save it in Redis
		newToken := uuid.NewString()
		rdb.Set(database.DbContext, newToken, userID, time.Hour)

		// Set the new session token in the Gin context
		c.Set("session_token", newToken)

		customWriter := &CustomResponseWriter{ResponseWriter: c.Writer, newSessionToken: newToken}
		c.Writer = customWriter

		c.Next()

		// Add the new session token to the JSON response
		customWriter.jsonResponse.SessionToken = customWriter.newSessionToken

		// Marshal the modified JSON response
		modifiedBody, err := json.Marshal(customWriter.jsonResponse)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			return
		}

		// Write the modified JSON response to the client
		customWriter.ResponseWriter.Write(modifiedBody)
	}
}
