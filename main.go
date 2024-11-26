package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JovidYnwa/microCmp/api"
	"github.com/JovidYnwa/microCmp/db"
	"github.com/JovidYnwa/microCmp/worker"
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
		Name:     os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
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

	var (
		companyStore         = db.NewPgCompanyStore(pgClient)
		companyHandler       = api.NewCompanyHandler(companyStore)
		companyFilterSotre   = db.NewOracleMainScreenStore(oracleClient)
		companyFilterHandler = api.NewCompanyFilterHandler(companyFilterSotre)

		dwhWorkerStore     = db.NewDwhWorkerStore(oracleClient)
		companyWorkerStore = db.NewWorkerStore(pgClient)
	)

	router := mux.NewRouter()
	router.HandleFunc("/filter/trpls", companyFilterHandler.HandleListTrpls)
	router.HandleFunc("/filter/regions", companyFilterHandler.HandleRgionsrpls)
	router.HandleFunc("/filter/subs/status", companyFilterHandler.HandleSubscriberStatus)
	router.HandleFunc("/filter/servs", companyFilterHandler.HandleServList)
	router.HandleFunc("/filter/sim/types", companyFilterHandler.HandleSimStatus)
	router.HandleFunc("/filter/device/types", companyFilterHandler.HandleDivceTypes)

	router.HandleFunc("/prize/list", companyFilterHandler.HandlePrizeList)
	router.HandleFunc("/action/list", companyFilterHandler.HandleActionCmp)

	router.HandleFunc("/company-type", companyHandler.HandleGetCompanies)
	router.HandleFunc("/company", companyHandler.HandleCreateCompany) //Post
	router.HandleFunc("/companies/{type_id:[0-9]+}", companyHandler.HandleGetCompany)

	router.HandleFunc("/company/{id:[0-9]+}", companyHandler.HandleGetCompanyDetail)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	//Wokers
	work := worker.NewCmpWoker("Cheaking unprocced comt", 20*time.Second, companyWorkerStore, dwhWorkerStore)
	go work.Start()

	log.Println("Json API server running on port: ", 3001)
	http.ListenAndServe(":3001", handler)
}
