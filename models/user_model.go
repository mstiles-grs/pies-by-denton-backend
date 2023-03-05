package models

import (
    // "database/sql"
    "log"
    "github.com/mstiles-grs/pies-by-denton-backend/database"
)


type User struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email  string `json:"email"`
    UserName  string `json:"user_name"`
    Password  string `json:"password"`
}

// CreateUserInDB creates a new user record in the database
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

