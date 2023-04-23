package ingredient

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/mstiles-grs/pies-by-denton-backend/database"
)

// This is what the obj of to store looks like

type Ingredient struct {
	UserID              string `json:"userID,omitempty" bson:"userID,omitempty"`
	IngredientName      string `json:"ingredientName"`
	IngredientType      string `json:"ingredientType"`
	StartingAmount      int    `json:"startingAmount"`
	StorageMeasurement  string `json:"storageMeasurement"`
	UnitOfMeasurement   string `json:"unitOfMeasurement"`
}

//This will be how I have to parse most of my Axios Requests
type QuickAddRequest struct {
	QuickAddIngredient Ingredient `json:"quickAddIngredient"`
}

// c* is for HTTP Requests/Data/Responses It also holds my userID
func QuickAdd(c *gin.Context) {

	var quickAddRequest QuickAddRequest

	// Parse the JSON object from the request body
	if err := c.ShouldBindJSON(&quickAddRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	quickAddRequest.QuickAddIngredient.UserID = userID.(string)

	// Access the "ingredients" collection using the existing connection
	collection := database.Client.Database("pies_by_denton").Collection("ingredients")

	// Insert the parsed data into the collection
	res, err := collection.InsertOne(database.DbContext, quickAddRequest.QuickAddIngredient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert the ingredient"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": res})
}