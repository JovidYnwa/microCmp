package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore1(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		db: db,
	}
}

type DatabaseConfig struct {
	Type     string
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=test host=db sslmode=disable port=5432"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func ConnectToPostgreSQL(config DatabaseConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", config.Type, config.User, config.Password, config.Host, config.Port, config.Name)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetConnMaxIdleTime(time.Minute * 5)

	// Verify the connection with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging Oracle: %w", err)
	}
	return db, nil
}

func ConnectToOracleGoOra(config DatabaseConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		config.Type,
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to Oracle: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetConnMaxIdleTime(time.Minute * 5)

	// Verify the connection with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging Oracle: %w", err)
	}
	return db, nil
}
