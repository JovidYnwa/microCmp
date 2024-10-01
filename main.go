package main

import (
	"log"
	"net/http"

	"github.com/JovidYnwa/microCmp/api"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/test1", api.HandleTestFunc1)

	log.Println("Json API server running on port: ", 3001)
	http.ListenAndServe(":3001", router)

	// store, err := db.NewPostgresStore()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := store.Init(); err != nil {
	// 	log.Fatalf("table was not created, %v", err)
	// }

	// server := api.NewAPIServer(":3001", store)
	// server.Run()
}
