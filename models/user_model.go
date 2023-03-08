package models

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/mstiles-grs/pies-by-denton-backend/database"
	"log"
	"time"
)

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUserInDB(user User) error {
	// get a connection to the database
	db := database.DB

	// prepare the INSERT statement
	statement, err := db.Prepare("INSERT INTO users (first_name, last_name, email, user_name, password) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	defer statement.Close()

	// execute the INSERT statement
	_, err = statement.Exec(user.FirstName, user.LastName, user.Email, user.UserName, user.Password)
	if err != nil {
		return err
	}

	log.Printf("User created: %s", user.UserName)
	return nil
}

func Login(email string, password string) (bool, string, error) {
	// Get a connection to the database
	db := database.DB

	// Query the database for a user with the given email and password
	var exists bool
	var userID int
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 AND password = $2), id FROM users WHERE email = $1 AND password = $2", email, password).Scan(&exists, &userID)

	if err != nil {
		return false, "", err
	}

	if exists {
		// Generate a unique session token
		sessionToken := uuid.New().String()

		// Set the session token and user ID in Redis
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		// Create a new context
		ctx := context.Background()

		err := rdb.Set(ctx, sessionToken, userID, time.Hour).Err()
		if err != nil {
			return false, "", err
		}

		return true, sessionToken, nil
	}

	return false, "", nil
}
