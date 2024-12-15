package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JovidYnwa/microCmp/api"
	"github.com/JovidYnwa/microCmp/db"
	"github.com/JovidYnwa/microCmp/internal/kafka"
	"github.com/JovidYnwa/microCmp/worker"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	// Import your local kafka package
	_ "github.com/sijms/go-ora/v2"
)

func main() {

	//TODO check pagination, job for updating statistic, sending notification function

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

	kafkaproducer := kafka.NewProducerFromEnv()
	defer kafkaproducer.Close()

	var (
		companyPgStore  = db.NewPgCompanyStore(pgClient)
		companyDwhStore = db.NewDwhWorkerStore(oracleClient)
		companyHandler  = api.NewCompanyHandler(companyPgStore, companyDwhStore)

		companyFilterSotre   = db.NewOracleMainScreenStore(oracleClient)
		companyFilterHandler = api.NewCompanyFilterHandler(companyFilterSotre)

		dwhWorkerStore     = db.NewDwhWorkerStore(oracleClient)
		companyWorkerStore = db.NewWorkerStore(pgClient)
	)

	cmpWorker := worker.NewCmpWoker(
		"Worker for setting new iteration for cmp",
		20*time.Hour,
		companyWorkerStore,
		dwhWorkerStore,
		worker.SetCmpIteration(companyWorkerStore),
	)

	cmpNotifierWorker := worker.NewCmpWoker(
		"Worker to send notification to subs who did not receive twice",
		20*time.Hour,
		companyWorkerStore,
		dwhWorkerStore,
		worker.CmpNotifier(companyWorkerStore, dwhWorkerStore),
	)

	cmpUpdateWorker := worker.NewCmpWoker(
		"Worker to update statistics",
		10*time.Minute,
		companyWorkerStore,
		dwhWorkerStore,
		worker.CmpStatisticUpdater(companyWorkerStore, dwhWorkerStore),
	)

	//Worker to update statistic

	go cmpWorker.Start()
	go cmpNotifierWorker.Start()
	go cmpUpdateWorker.Start()

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
	// work := worker.NewCmpWoker("Cheaking unprocced comt", 20*time.Second, companyWorkerStore, dwhWorkerStore)
	// go work.Start()
	// After successfully creating the company, send a message to Kafka

	log.Println("Json API server running on port: ", 3001)
	http.ListenAndServe(":3001", handler)

}
