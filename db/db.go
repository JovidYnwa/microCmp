package db

import (
	"database/sql"
	"fmt"
	"log"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=test host=db sslmode=disable port=5432"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Println("we are fucked")
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}
