package middleware

import (
	"context"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"

)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Println("AuthMiddleware called for URL path:", c.Request.URL.Path)
        // Skip authentication check for login and create user routes
        if c.Request.URL.Path == "/login/user" || c.Request.URL.Path == "/create/user" {
            c.Next()
            return
        }

        log.Println(c.Request.Method, c.Request.URL.Path)
		log.Println("Headers:", c.Request.Header)


        sessionID, err := c.Cookie("session_token")
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        // Get the user ID associated with the session token
        rdb := redis.NewClient(&redis.Options{
            Addr:     "localhost:6379",
            Password: "",
            DB:       0,
        })
        ctx := context.Background()
        userID, err := rdb.Get(ctx, sessionID).Result()
        if err != nil {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        // Set the user ID in the request context for downstream handlers to use
        c.Set("userID", userID)

        c.Next()
    }
}