package controllers

import (
	"net/http"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mstiles-grs/pies-by-denton-backend/models"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/mstiles-grs/pies-by-denton-backend/models/recipe"
    "github.com/mstiles-grs/pies-by-denton-backend/models/ingredient"
)

type UserController struct{}

// CreateUser is the handler function for the POST /create/user route
func (u *UserController) CreateUser(c *gin.Context, client *mongo.Client, ctx context.Context) {
	// Parse the JSON request body into a User struct
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Create the user in the database
	if err := models.CreateUserInDB(client, ctx, newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return a success message and the new user object
	c.JSON(http.StatusCreated, gin.H{"message": "user created", "user": newUser})
}

func (u *UserController) LoginUser(c *gin.Context, client *mongo.Client, ctx context.Context) {
	// Parse the JSON request body into a LoginUser struct
	var loginUser models.LoginUser
	if err := c.BindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Call the login function on the user model
	success, sessionToken, err := models.Login(client, ctx, loginUser.Email, loginUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Return a failure response if login was unsuccessful
	if !success {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Set the session token in a cookie
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600,
	}
	http.SetCookie(c.Writer, cookie)

	// Return a success response with the session token in the body
	c.JSON(http.StatusOK, gin.H{"message": "login successful", "session_token": sessionToken})
}

func (u *UserController) Dashboard(c *gin.Context, client *mongo.Client, ctx context.Context) {
	
    userID, err := primitive.ObjectIDFromHex(c.GetString("userID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    recipes, err := recipe.GetUserRecipes(client, ctx, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user recipes"})
        return
    }

    ingredients, err := ingredient.GetUserIngredients(client, ctx, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user ingredients"})
        return
    }

    if len(recipes) == 0 && len(ingredients) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": "Recipes and ingredients are empty"})
    } else {
        c.JSON(http.StatusOK, gin.H{"recipes": recipes, "ingredients": ingredients})
    }
}
