package main

import (
	"log"
	"net/http"

	"github.com/JovidYnwa/microCmp/api"
	"github.com/JovidYnwa/microCmp/db"
	"github.com/gorilla/mux"

	_ "github.com/sijms/go-ora/v2" // Import the go-ora driver
)

func main() {

	dwhConfigs := db.DatabaseConfig{
		Type:     "oracle",
		Name:     "2345",
		Host:     "134560.832454.312343.222435",
		Port:     "1521",
		User:     "фва",
		Password: "ыап$ыавп",
	}

	oracleClient, err := db.ConnectToOracleGoOra(dwhConfigs)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer func() {
		err := oracleClient.Close()
		if err != nil {
			log.Fatalf("Can't close go-ora connection: %s", err)
		}
	}()

	companySotre := db.NewOracleMainScreenStore(oracleClient)
	companyHandler := api.NewCompanyHandler(companySotre)

	router := mux.NewRouter()
	router.HandleFunc("/filter/trpls", companyHandler.HandleListTrpls)
	router.HandleFunc("/filter/regions", companyHandler.HandleRgionsrpls)
	router.HandleFunc("/filter/subs/status", companyHandler.HandleSubscriberStatus)

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
