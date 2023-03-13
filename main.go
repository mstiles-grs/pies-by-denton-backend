package main

import (
    "log"
    "os"
    "io"
    "net/http"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/mstiles-grs/pies-by-denton-backend/controllers"
    "github.com/mstiles-grs/pies-by-denton-backend/database"
	"github.com/mstiles-grs/pies-by-denton-backend/middleware"
	"github.com/mstiles-grs/pies-by-denton-backend/models"
	"bytes"
    "encoding/json"

)


func main() {


    // Initialize the database and defer closing the connection
    database.InitDB()
    defer database.DB.Close()

    // Open the log file for writing
    logFile, err := os.OpenFile("request.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer logFile.Close()

    // Use a multi-writer to write logs to both the console and the log file
    logWriter := io.MultiWriter(logFile, os.Stdout)
    log.SetOutput(logWriter)

    router := gin.Default()

    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://127.0.0.1:5173"}
    config.AllowCredentials = true
    config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}
    router.Use(cors.New(config))

    // Use the AuthMiddleware for all routes
    router.Use(middleware.AuthMiddleware())


    //Using a Blank Struct in UserController to handle an format


	userController := &controllers.UserController{}

	router.POST("/create/user", func(c *gin.Context) {
		// Read the request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		// Restore the request body for downstream middleware/handlers
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Log the incoming data
		log.Printf("Request: %s\n", string(body))

		// Parse the request body into a User struct
		var newUser models.User
		if err := json.Unmarshal(body, &newUser); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Call the controller function
		userController.CreateUser(c)
	})

	router.POST("/login/user", func(c *gin.Context) {
        // Read the request body
        body, err := io.ReadAll(c.Request.Body)
        if err != nil {
            log.Println(err)
            c.AbortWithStatus(http.StatusBadRequest)
            return
        }
        // Restore the request body for downstream middleware/handlers
        c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

        // Log the incoming data
        log.Printf("Request: %s\n", string(body))

        // Parse the request body into a LoginUser struct
        var loginUser models.LoginUser
        if err := json.Unmarshal(body, &loginUser); err != nil {
            log.Println(err)
            c.AbortWithStatus(http.StatusBadRequest)
            return
        }

        // Call the controller function
        userController.LoginUser(c)
    })

    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }



    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

