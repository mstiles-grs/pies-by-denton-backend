package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

)

var Client *mongo.Client
var DbContext context.Context

func InitDB() error {
	connectionString := "mongodb://localhost:27017"

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("error pinging MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")

	Client = client
	DbContext = context.Background()

	return nil
}

func ValidateUserCredentials(email, password string) (bool, string) {
	usersCollection := Client.Database("pies_by_denton").Collection("users")

	filter := bson.M{"email": email}
	var user struct {
		ID       string `bson:"_id"`
		Email string `bson:"email"`
		Password string `bson:"password"`
	}

	err := usersCollection.FindOne(DbContext, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, ""
		}
		fmt.Printf("Error querying user: %v\n", err)
		return false, ""
	}

	// Compare the input password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// The password doesn't match
		return false, ""
	}

	return true, user.ID
}