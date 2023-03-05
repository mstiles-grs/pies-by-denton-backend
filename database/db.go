package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// DB is the global database connection variable
var DB *sql.DB

func InitDB() error {
    // get the database url from the environment variable
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL environment variable is not set")
    }

    // append sslmode=disable to the database URL to disable SSL
    dbURL += "?sslmode=disable"

    // open a connection to the database
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        return fmt.Errorf("error opening database connection: %v", err)
    }

    // test the connection
    err = db.Ping()
    if err != nil {
        return fmt.Errorf("error connecting to database: %v", err)
    }

    fmt.Println("Successfully connected to database")

    // set the global database connection variable
    DB = db

    return nil
}
