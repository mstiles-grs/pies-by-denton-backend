package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// initialize the gin router
	router := gin.Default()

	// set up a GET route to test the server
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
	})

	// set up a POST route to handle creating a new user
	router.POST("/users", createUser)

	// start the server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Error starting server: ", err.Error())
	}
}

// createUser is the handler function for the POST /users route
func createUser(c *gin.Context) {
	// parse the JSON request body into a User struct
	var newUser User
	err := c.BindJSON(&newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create the user in the database
	err = createUserInDB(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return a success message and the new user object
	c.JSON(http.StatusCreated, newUser)
}

// createUserInDB is a function that creates a new user in the database
func createUserInDB(newUser User) error {
	// here you would add code to create the user in your PostgreSQL database
	// for example, using the "database/sql" package

	// for now, just print out the new user to the console
	fmt.Println("Created new user:", newUser)
	return nil
}