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

func InitDB() {
	// get the database url from the environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening database connection: ", err.Error())
	}

	// test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database: ", err.Error())
	}

	fmt.Println("Successfully connected to database")

	// set the global database connection variable
	DB = db
}