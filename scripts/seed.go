package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Get database connection details from environment variables
	dbType := os.Getenv("DB_TYPE")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	// Construct the connection string
	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		dbType, dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	if err := seedDB(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
	fmt.Println("Database seeded successfully!")
}

func seedDB(db *sql.DB) error {
	// Insert seed data
	_, err := db.Exec(`
		INSERT INTO company (NAME, company_launched, subscriber_count, efficiency) 
		VALUES 
		('Mk <<No Pay>>', 5, 500, 50.2),
		('Mk <<Must Pay>>', 6, 501, 50.2),
		('Mk <<Suffer>>', 5, 500, 50.2);
	`)
	return err
}
