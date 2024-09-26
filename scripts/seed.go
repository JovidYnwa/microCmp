package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
)

// should get from env
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "test"
	dbname   = "postgres"
)

func main() {
	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()

	//Connect to the database
	cnnstr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", cnnstr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to the db!")

	if *seed {
		err = seedDB(db)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("db seeded successfully")
	}
}


func seedDB(db *sql.DB) error {
	// Insert seed data
	_, err := db.Exec(`
		INSERT INTO TEST_TABLE (NAME, AGE, EMAIL) 
		VALUES 
		('John Doe', 30, 'john.doe@example.com'),
		('Jane Smith', 25, 'jane.smith@example.com'),
		('Bob Johnson', 45, 'bob.johnson@example.com')
		ON CONFLICT (NAME) DO NOTHING;
	`)
	return err
}