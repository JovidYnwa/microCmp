package main

import (
	"log"

	"github.com/JovidYnwa/microCmp/api"
	"github.com/JovidYnwa/microCmp/db"
)

func main() {
	log.Println("Json API server running on port: ")
	store, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatalf("table was not created, %v", err)
	}

	server := api.NewAPIServer(":3001", store)
	server.Run()
}
