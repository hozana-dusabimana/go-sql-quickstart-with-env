// Package main provides a PostgreSQL database connection example.
//
// This application demonstrates:
// - Connecting to a PostgreSQL database using pgx
// - Loading configuration from environment variables using Viper
// - Creating tables with schema constraints
// - Inserting data with duplicate-key conflict handling
package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/spf13/viper"
	"github.com/jackc/pgx/v5"
	//use godotenv to load .env file
	// "github.com/joho/godotenv"
	// "os"
)

// main is the entry point of the application.
// It performs the following steps:
// 1. Loads configuration from .env file using Viper
// 2. Connects to PostgreSQL database
// 3. Creates a users table if it doesn't exist
// 4. Inserts sample user records with conflict handling
// 5. Displays results and configuration values
func main() {

	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// connStr := os.Getenv("CONN_STR")

	// use viper to load .env file
	// Viper is used for configuration management, providing flexibility
	// to load from environment variables, config files, and more

	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // read in environment variables that match
	viper.Set("Developer", "Hozana")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Retrieve the connection string from configuration
	connStr := viper.GetString("CONN_STR")

	// Connect to PostgreSQL database using pgx
	// context.Background() is used as the base context for the connection
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	// Ensure the connection is properly closed when the function returns
	defer conn.Close(context.Background())

	// Query the current database time to verify connection
	var now time.Time
	err = conn.QueryRow(context.Background(), "SELECT NOW()").Scan(&now)
	if err != nil {
		log.Fatal("QueryRow failed:", err)
	}

	// SQL statement to create the users table
	// IF NOT EXISTS ensures idempotency - the table is only created if it doesn't exist
	// UNIQUE constraints on username and email prevent duplicate entries
	// created_at automatically records when each record is inserted
	tablecreate := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Execute the table creation statement
	if _, err := conn.Exec(context.Background(), tablecreate); err != nil {
		log.Fatal("Table creation failed:", err)
	}

	fmt.Println("Table 'users' created or already exists.")

	// Sample user data to insert
	// Note: The third user has the same username as the first, which will test conflict handling
	users := []map[string]string{
		{"username": "alice", "email": "alice@example.com"},
		{"username": "bob", "email": "bob@example.com"},
		{"username": "alice", "email": "alice@example.com"}, // duplicate username
	}

	// SQL statement for inserting users with conflict resolution
	// ON CONFLICT (username) DO NOTHING silently ignores duplicate username insertions
	// This prevents the application from crashing on duplicate entries
	// $1 and $2 are parameterized placeholders for username and email respectively
	addUserSql := `INSERT INTO users (username, email)
	               VALUES ($1, $2)
	               ON CONFLICT (username) DO NOTHING;`

	// Iterate through users and attempt to insert each one
	// Errors are logged and the loop continues, allowing partial success
	for _, user := range users {
		_, err := conn.Exec(context.Background(), addUserSql, user["username"], user["email"])
		if err != nil {
			log.Printf("Failed to insert user %s: %v", user["username"], err)
			continue
		}
		fmt.Printf("User %s inserted successfully\n", user["username"])
	}

	// Display the current database time and configuration
	fmt.Println("Current time:", now)
	fmt.Println("Developer:", viper.GetString("Developer"))
}
