package models

import (
	"context"
	"time"
	"log"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Email     string             `bson:"email" json:"email"`
	UserName  string             `bson:"user_name" json:"user_name"`
	Password  string             `bson:"password" json:"password"`
}

type LoginUser struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

// TableName sets the table name for the User model


func CreateUserInDB(client *mongo.Client, ctx context.Context, user User) error {

	usersCollection := client.Database("pies_by_denton").Collection("users")
	log.Printf("Password before hashing: %s", user.Password)

	// Hash the user's password
	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	log.Printf("Password after hashing: %s", user.Password)

	_, err = usersCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}


func Login(client *mongo.Client, ctx context.Context, email string, password string) (bool, string, error) {

	usersCollection := client.Database("pies_by_denton").Collection("users")



	var user User
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)


	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, "", nil
		}
		return false, "", err
	}
	// Compare the hashed password from the database with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Log the error message

		return false, "", nil // Passwords do not match
	}


	// Generate a unique session token
	sessionToken := uuid.New().String()


	// Set the session token and user ID in Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Error pinging Redis: %v", err)
	} else {
		log.Printf("Redis ping result: %s", pong)
	}

	err = rdb.Set(ctx, sessionToken, user.ID.Hex(), time.Hour).Err()
	if err != nil {
		return false, "", err
	}



	return true, sessionToken, nil
}