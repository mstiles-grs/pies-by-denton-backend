package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/mstiles-grs/pies-by-denton-backend/controllers"
	"github.com/mstiles-grs/pies-by-denton-backend/db"

)

func main() {

	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.DB.Close()


    // initialize the gin router
    router := gin.Default()

    // create a new UserController
    userController := &controllers.UserController{}

    // set up a POST route to handle creating a new user
    router.POST("/create/user", userController.CreateUser)

    // start the server
    err := router.Run(":8080")
    if err != nil {
        log.Fatal("Error starting server: ", err.Error())
    }
}