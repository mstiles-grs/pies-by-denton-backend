package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/mstiles-grs/pies-by-denton-backend/database"
	"time"
	"strings"
    "errors"
)

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type CheckRequest struct {
    Token    string `json:"session_token"`
}


func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// If the route is /create/user, ignore authentication
		if path == "/create/user" {
			c.Next()
			return
		}

		if path == "/login/user" {
			var loginRequest LoginRequest
            if err := c.BindJSON(&loginRequest); err != nil {
                c.AbortWithStatusJSON(400, gin.H{"error": "Invalid request body"})
                return
            }
            email := loginRequest.Email
            password := loginRequest.Password


			// Check the user credentials against the MongoDB database
			isValid, userID := validateCredentials(email, password)

			if isValid {
				// Generate a UUID and associate it with the user ID in Redis
				token := uuid.NewString()
				err := setTokenInRedis(userID, token)
				if err == nil {
					c.JSON(200, gin.H{
						"session_token": token,
					})
					return
				}
			}
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		CheckSessionMiddleware()(c)
	}
}

func validateCredentials(email, password string) (bool, string) {

	return database.ValidateUserCredentials(email, password)

}

func setTokenInRedis(userID, token string) error {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Set the token associated with the userID
	err := rdb.Set(database.DbContext, token, userID, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func CheckSessionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get the session token from the Authorization header
        sessionToken := c.Request.Header.Get("Authorization")

        // Check if the session token is present and valid
        userID, err := validateSessionToken(sessionToken)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid session token"})
            return
        }

        // Attach the user ID to the request context
        c.Set("userID", userID)

        // Call the next middleware function or route handler function
        c.Next()
    }
}

func validateSessionToken(sessionToken string) (string, error) {

    // Remove the "Bearer " prefix from the sessionToken
    if len(sessionToken) > 7 && strings.HasPrefix(sessionToken, "Bearer ") {
        sessionToken = sessionToken[7:]
    } else {
        return "", errors.New("Invalid session token format")
    }

    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    // Get the userID associated with the session token from Redis
    userID, err := rdb.Get(database.DbContext, sessionToken).Result()

    if err != nil {
        return "", err
    }

    return userID, nil
}



