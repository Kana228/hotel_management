package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

func Connect() {
	connStr := "host=localhost port=5434 user=postgres dbname=mydb sslmode=disable"
	var err error

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database is not reachable: %v", err)
	}

	fmt.Println("Successfully connected to the database.")
}

// Close the database connection when the app stops
func Close() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed.")
	}
}
