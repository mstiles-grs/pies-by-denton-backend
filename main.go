package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mstiles-grs/pies-by-denton-backend/user"
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
	// router.Use(middleware.RefreshSessionMiddleware())
	router.POST("/login/user", func(c *gin.Context) {
		c.Next()
	})
	// router.POST("/create/user", func(c *gin.Context) {
	// 	users.CreateUser(c, database.Client, database.DbContext)
	// })

	router.GET("/user/dashboard", middleware.SessionMiddleware(), func(c *gin.Context) {
		user.Dashboard(c)
	})

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}