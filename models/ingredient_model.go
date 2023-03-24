package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Ingredient struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	measurment string             `bson:"measurment" json:"measurment"`
	price int32             `bson:"price" json:"price"`
	quantity  string             `bson:"quantity" json:"quantity"`
	common     bool             `bson:"common" json:"common"`
	name  string             `bson:"name" json:"name"`
	base_substance  string             `bson:"base_substance" json:"base_substance"`
	entered_by 	primitive.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
}

func GetUserIngredients(client *mongo.Client, ctx context.Context, userID primitive.ObjectID) ([]Ingredient, error) {
    ingredientsCollection := client.Database("pies_by_denton").Collection("ingredients")

    cursor, err := ingredientsCollection.Find(ctx, bson.M{"user_id": userID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var ingredients []Ingredient
    if err = cursor.All(ctx, &ingredients); err != nil {
        return nil, err
    }

    return ingredients, nil
}