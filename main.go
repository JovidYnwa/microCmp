package main

import (
	"log"
	"net/http"
	"os"

	"github.com/JovidYnwa/microCmp/api"
	"github.com/JovidYnwa/microCmp/db"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	_ "github.com/sijms/go-ora/v2"
)

func main() {

	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dwhConfigs := db.DatabaseConfig{
		Type:     os.Getenv("DWH_DB_TYPE"),
		Name:     os.Getenv("DWH_DB_NAME"),
		Host:     os.Getenv("DWH_DB_HOST"),
		Port:     os.Getenv("DWH_DB_PORT"),
		User:     os.Getenv("DWH_DB_USER"),
		Password: os.Getenv("DWH_DB_PASSWORD"),
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

	pgConfigs := db.DatabaseConfig{
		Type:     os.Getenv("DB_TYPE"),
		Name:     os.Getenv("DB_NAME"),     //"postgres",
		Host:     os.Getenv("DB_HOST"),     //"db",
		Port:     os.Getenv("DB_PORT"),     //"5432",
		User:     os.Getenv("DB_USER"),     //"postgres",
		Password: os.Getenv("DB_PASSWORD"), //"test",
	}

	pgClient, err := db.ConnectToPostgreSQL(pgConfigs)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := pgClient.Close()
		if err != nil {
			log.Fatalf("Can't close pg connection: %s", err)
		}
	}()

	companyStore := db.NewPgCompanyStore(pgClient)
	companyHandler := api.NewCompanyHandler(companyStore)

	companyFilterSotre := db.NewOracleMainScreenStore(oracleClient)
	companyFilterHandler := api.NewCompanyFilterHandler(companyFilterSotre)

	router := mux.NewRouter()
	router.HandleFunc("/filter/trpls", companyFilterHandler.HandleListTrpls)
	router.HandleFunc("/filter/regions", companyFilterHandler.HandleRgionsrpls)
	router.HandleFunc("/filter/subs/status", companyFilterHandler.HandleSubscriberStatus)

	router.HandleFunc("/companies", companyHandler.HandleGetCompanies)
	router.HandleFunc("/company", companyHandler.HandleCreateCompany)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Println("Json API server running on port: ", 3001)
	http.ListenAndServe(":3001", handler)

	// store, err := db.NewPostgresStore()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := store.Init(); err != nil {
	// 	log.Fatalf("table was not created, %v", err)
	// }

	// server := api.NewCompanyHandler(":3001", store)
	// server.Run()
}
