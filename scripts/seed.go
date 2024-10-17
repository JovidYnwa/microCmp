package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

	// if err := seedCompany(db); err != nil {
	// 	log.Fatalf("Failed to seed company: %v", err)
	// }
	// fmt.Println("Company table seeded successfully!")

	if err := seedCompanyInfo(db); err != nil {
		log.Fatalf("Failed to seed company_info: %v", err)
	}
	fmt.Println("Company info seeded successfully!")

	if err := seedCompanyRepetition(db); err != nil {
		log.Fatalf("Failed to seed company_repetition: %v", err)
	}
	fmt.Println("Company repetition seeded successfully!")
}

func seedCompany(db *sql.DB) error {
	query := `
        INSERT INTO company (
            cmp_name,
            cmp_description,
            navi_user,
            query_id,
            start_time,
            duration,
            repetition
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	// Sample company names and users
	companies := []struct {
		name     string
		desc     string
		naviUser string
	}{
		{"Mk <<Tech Solutions>>", "Описание  1", "Mr Seedr"},
		{"Mk <<Digital Innovation>>", "Описание 2", "Mr Seedr"},
		{"Mk <<Smart Systems>>", "Описание 3", "Mr Seedr"},
		{"Mk <<Cloud Computing>>", "Описание 4", "Mr Seedr"},
		{"Mk <<Data Analytics>>", "Описание 5", "Mr Seedr"},
	}

	baseTime := time.Now()
	for i, company := range companies {
		_, err := db.Exec(query,
			company.name,
			company.desc,
			company.naviUser,
			fmt.Sprintf("QID-%d", i+1), // Unique query ID
			baseTime.Add(time.Duration(i*24)*time.Hour), // Staggered start times
			10+i, // Different durations
			5+i,  // Different repetitions
		)
		if err != nil {
			return fmt.Errorf("failed to insert company record: %v", err)
		}
	}
	return nil
}

func seedCompanyInfo(db *sql.DB) error {
	query := `
        INSERT INTO company_info (
            company_id,
            cmp_filter,
            sms_data,
            action_data
        ) VALUES (
            $1, $2, $3, $4
        )
    `

	// Sample data for JSON fields
	phoneTypes := []types.BaseFilter{{ID: 1, Name: "Mobile"}, {ID: 2, Name: "Home"}}
	trpls := []types.BaseFilter{{ID: 1, Name: "Plan A"}, {ID: 2, Name: "Plan B"}}
	balanceLimits := types.BalanceLimit{Start: 50.0, End: 200.0}
	subscriberStatuses := []types.BaseFilter{{ID: 1, Name: "Active"}, {ID: 2, Name: "Suspended"}}
	deviceType := 1
	packSpent := types.TrafficSpent{Min: 100, Sms: 50, MB: 1024}
	arpuLimits := types.ARPULimit{Start: 10.0, End: 100.0}
	regions := []types.BaseFilter{{ID: 1, Name: "Dushanbe"}, {ID: 2, Name: "Khujand"}}
	simDate := time.Now().AddDate(0, -6, 0) // Sim date 6 months ago
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

	// Serialize the JSON fields
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

	// Get the company IDs
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

	// Insert the data into the company_info table
	for _, companyID := range companyIDs {
		// Check if the record already exists
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM company_info WHERE company_id = $1)", companyID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check existence of company_info record: %v", err)
		}

		if exists {
			fmt.Printf("Skipping existing company_info record for company_id: %d\n", companyID)
			continue // Skip this record if it already exists
		}

		_, err = db.Exec(query, companyID, cmpFilterJSON, smsDataJSON, actionDataJSON)
		if err != nil {
			return fmt.Errorf("failed to insert company_info record: %v", err)
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
            start_date,
            end_date
        ) VALUES ($1, $2, $3, $4, $5)
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
		endDate := startDate.AddDate(0, 1, 0) // One month later

		_, err := db.Exec(query,
			companyID,
			efficiency,
			subAmount,
			startDate,
			endDate,
		)
		if err != nil {
			return fmt.Errorf("failed to insert company_repetion record: %v", err)
		}
	}

	return nil
}
