package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mstiles-grs/pies-by-denton-backend/controllers"
	"github.com/mstiles-grs/pies-by-denton-backend/database"
	"github.com/mstiles-grs/pies-by-denton-backend/middleware"
	"time"
)

func main() {
	// Initialize the database and defer closing the connection
	err := database.InitDB()
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173"
		},
		MaxAge: 12 * time.Hour,
	}))

	// Use the AuthMiddleware for all routes
	router.Use(middleware.AuthMiddleware())

	// Using a Blank Struct in UserController to handle an forat
	userController := &controllers.UserController{}

	router.POST("/create/user", func(c *gin.Context) {
		// Call the controller function with the database connection from the context
		userController.CreateUser(c, database.Client, database.DbContext)
	})

	router.POST("/login/user", func(c *gin.Context) {
		// Call the controller function with the database connection from the context
		userController.LoginUser(c, database.Client, database.DbContext)
	})

	router.GET("/user/dashboard", func(c *gin.Context) {
		userController.Dashboard(c, database.Client, database.DbContext)
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
