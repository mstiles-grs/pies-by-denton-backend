package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/mstiles-grs/pies-by-denton-backend/models"
)

type UserController struct{}

// CreateUser is the handler function for the POST /create/user route
func (u *UserController) CreateUser(c *gin.Context) {
    // Parse the JSON request body into a User struct
    var newUser models.User
    if err := c.BindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    // Create the user in the database
    if err := models.CreateUserInDB(newUser); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
        return
    }

    // Return a success message and the new user object
    c.JSON(http.StatusCreated, gin.H{"message": "user created", "user": newUser})
}