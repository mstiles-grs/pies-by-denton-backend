package controllers


import (
    "github.com/gin-gonic/gin"
   "~/go/src/github.com/mstiles-grs/pies-by-denton-backend/models"
)

// CreateUserHandler handles requests to create a new user.
func CreateUserHandler(c *gin.Context) {
	// Parse request body as JSON
	var newUser user_model.User
	if err := json.NewDecoder(c.Request.Body).Decode(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Save new user to database
	if err := newUser.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}