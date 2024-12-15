package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/JovidYnwa/microCmp/types"
	_ "github.com/lib/pq"
)

func main() {
	dbType := os.Getenv("DB_TYPE")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		dbType, dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	if err := seedCompanyType(db); err != nil {
		log.Fatalf("Failed to seed company: %v", err)
	}
	fmt.Println("Company table seeded successfully!")

	if err := seedCompany(db); err != nil {
		log.Fatalf("Failed to seed company_info: %v", err)
	}
	fmt.Println("Company info seeded successfully!")

	if err := seedCompanyRepetition(db); err != nil {
		log.Fatalf("Failed to seed company_repetition: %v", err)
	}
	fmt.Println("Company repetition seeded successfully!")
}

func seedCompanyType(db *sql.DB) error {
	query := `
        INSERT INTO company_type(
            cmp_name,
            navi_user
        ) VALUES ($1, $2)
    `

	// Sample company names and users
	companies := []struct {
		name     string
		naviUser string
	}{
		{"Mk <<No Pay", "Mr Seedr"},
		{"Mk <<Must Pay>>", "Mr Seedr"},
		{"Mk <<Suffer>>", "Mr Seedr"},
	}

	for _, company := range companies {
		_, err := db.Exec(query,
			company.name,
			company.naviUser,
		)
		if err != nil {
			return fmt.Errorf("failed to insert company record: %v", err)
		}
	}
	return nil
}

func seedCompany(db *sql.DB) error {
	query := `
        INSERT INTO company (
            company_type_id,
            cmp_billing_id,
            start_date,
            end_date,
            cmp_desc,
            cmp_filter,
            sms_data,
            action_data
        ) VALUES (
            $1, $2, $3, $4, $5::jsonb, $6::jsonb, $7::jsonb, $8::jsonb
        )
    `

	// Sample data for JSON fields (as provided)
	phoneTypes := []types.BaseFilter{{ID: 1, Name: "gold"}, {ID: 2, Name: "sliver"}}
	trpls := []types.BaseFilter{{ID: 1, Name: "salom 1"}, {ID: 2, Name: "salom 2"}}
	balanceLimits := types.BalanceLimit{Start: 50.0, End: 200.0}
	subscriberStatuses := []types.BaseFilter{{ID: 1, Name: "Active"}, {ID: 2, Name: "Suspended"}}
	deviceType := []types.BaseFilter{{ID: 1, Name: "android"}, {ID: 2, Name: "ios"}}
	packSpent := types.TrafficSpent{Min: 100, Sms: 50, MB: 1024}
	arpuLimits := types.ARPULimit{Start: 10.0, End: 100.0}
	regions := []types.BaseFilter{{ID: 1, Name: "Dushanbe"}, {ID: 2, Name: "Khujand"}}
	simDate := types.CustomTime{Time: time.Now().AddDate(0, -6, 0)} // Sim date 6 months ago
	services := []types.BaseFilter{{ID: 1, Name: "Service A"}, {ID: 2, Name: "Service B"}}
	wheelUsage := true

	cmpFilter := types.CompanyInfo{
		PhoneType:        phoneTypes,
		Trpl:             trpls,
		BalanceLimits:    balanceLimits,
		SubscriberStatus: subscriberStatuses,
		DeviceType:       deviceType,
		PackSpent:        packSpent,
		ARPULimits:       arpuLimits,
		Region:           regions,
		SimDate:          simDate,
		Service:          services,
		WheelUsage:       wheelUsage,
	}

	smsData := types.SmsBefore{
		SmsText: types.TextType{Ru: "Текст уведомления", Tj: "Матни хабар", Eng: "Reminder text"},
		SmsDay:  3,
		SmsTextRemid: types.TextType{
			Ru:  "Текст напоминания",
			Tj:  "Матни ёдраси",
			Eng: "Reminder message",
		},
	}

	actionData := types.CompanyAction{
		Action: types.BaseFilter{ID: 1, Name: "Send SMS"},
		Sms:    types.TextType{Ru: "Сообщение для действия", Tj: "Паёми барои амал", Eng: "Action message"},
		Prize:  types.BaseFilter{ID: 2, Name: "Free Data"},
	}

	// Data for cmp_desc
	cmpDesc := map[string]interface{}{
		"desc":        "Leading tech innovation company",
		"name":        "Tech Innovators Ltd.2",
		"repition":    3,
		"startTime":   "2024-10-17T10:00:00Z",
		"durationDay": 30,
	}

	// Serialize the JSON fields
	cmpDescJSON, err := json.Marshal(cmpDesc)
	if err != nil {
		return fmt.Errorf("failed to marshal cmp_desc: %v", err)
	}
	cmpFilterJSON, err := json.Marshal(cmpFilter)
	if err != nil {
		return fmt.Errorf("failed to marshal cmp_filter: %v", err)
	}
	smsDataJSON, err := json.Marshal(smsData)
	if err != nil {
		return fmt.Errorf("failed to marshal sms_data: %v", err)
	}
	actionDataJSON, err := json.Marshal(actionData)
	if err != nil {
		return fmt.Errorf("failed to marshal action_data: %v", err)
	}

	// Set the start and end dates
	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, 5) // End date set to 1 month after start

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Insert the data into the company table
	for i := 0; i < 5; i++ { // Insert 10 records
		// Generate a random company_type_id between 1 and 100
		companyTypeID := 1
		// Generate a random cmp_billing_id between 1000 and 9999
		cmpBillingID := rand.Intn(9000) + 1000

		_, err = db.Exec(query, companyTypeID, cmpBillingID, startDate, endDate, cmpDescJSON, cmpFilterJSON, smsDataJSON, actionDataJSON)
		if err != nil {
			return fmt.Errorf("failed to insert company record: %v", err)
		}
	}

	return nil
}

func seedCompanyRepetition(db *sql.DB) error {
	query := `
        INSERT INTO company_repetion (
            company_id,
            efficiency,
            sub_amount,
            start_date
            --end_date
        ) VALUES ($1, $2, $3, $4)
    `

	rows, err := db.Query("SELECT id FROM company ORDER BY id LIMIT 10")
	if err != nil {
		return fmt.Errorf("failed to get company IDs: %v", err)
	}
	defer rows.Close()

	var companyIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan company ID: %v", err)
		}
		companyIDs = append(companyIDs, id)
	}

	baseTime := time.Now()
	for i, companyID := range companyIDs {
		efficiency := 0.5 + float64(i)*0.05 // Efficiency from 0.5 to 0.95
		subAmount := 1000 + i*100           // Sub amount from 1000 to 1900
		startDate := baseTime.AddDate(0, 0, -i)
		//endDate := startDate.AddDate(0, 1, 0) // One month later

		_, err := db.Exec(query,
			companyID,
			efficiency,
			subAmount,
			startDate,
			//endDate,
		)
		if err != nil {
			return fmt.Errorf("failed to insert company_repetion record: %v", err)
		}
	}

	return nil
}
