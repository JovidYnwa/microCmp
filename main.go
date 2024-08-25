package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatalf("table was not created, %v", err)
	}

	server := NewAPIServer(":3001", store)
	server.Run()
}
