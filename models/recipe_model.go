package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

)

type Recipe struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	name                 string             `bson:"first_name" json:"first_name"`
	ingredient           string             `bson:"last_name" json:"last_name"`
	ppu                  int32              `bson:"email" json:"email"`
	created_by           primitive.ObjectID `bson:"_id,created_by" json:"_id,created_by"`
	ingredients_avliable bool               `bson:"ingredient_avliable" json:"ingredient_avliable"`
}

func GetUserRecipes(client *mongo.Client, ctx context.Context, userID primitive.ObjectID) ([]Recipe, error) {
    recipesCollection := client.Database("pies_by_denton").Collection("recipes")

    cursor, err := recipesCollection.Find(ctx, bson.M{"user_id": userID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var recipes []Recipe
    if err = cursor.All(ctx, &recipes); err != nil {
        return nil, err
    }

    return recipes, nil
}
