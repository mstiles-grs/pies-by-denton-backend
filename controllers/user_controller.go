
package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/mstiles-grs/pies-by-denton-backend/models"
)

type UserController struct {}

// createUser is the handler function for the POST /create/user route
func (u *UserController) CreateUser(c *gin.Context) {
    // parse the JSON request body into a User struct
    var newUser models.User
    err := c.BindJSON(&newUser)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // create the user in the database
    if err := models.CreateUserInDB(newUser); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // return a success message and the new user object
    c.JSON(http.StatusCreated, newUser)
}