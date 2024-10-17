package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"

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

    if err := seedCompany(db); err != nil {
        log.Fatalf("Failed to seed company: %v", err)
    }
    fmt.Println("Company table seeded successfully!")

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
        {"Mk <<Tech Solutions>>", "Описание  1","Mr Seedr"},
        {"Mk <<Digital Innovation>>", "Описание 2","Mr Seedr"},
        {"Mk <<Smart Systems>>", "Описание 3","Mr Seedr"},
        {"Mk <<Cloud Computing>>", "Описание 4","Mr Seedr"},
        {"Mk <<Data Analytics>>", "Описание 5","Mr Seedr"},
        {"Mk <<Cyber Security>>", "Описание 6","Mr Seedr"},
        {"Mk <<Web Services>>", "Описание 7","Mr Seedr"},
        {"Mk <<Mobile Solutions>>", "Описание 8","Mr Seedr"},
        {"Mk <<AI Systems>>", "Описание 9","Mr Seedr"},
        {"Mk <<IoT Platform>>", "Описание 10","Mr Seedr"},
    }

    baseTime := time.Now()
    for i, company := range companies {
        _, err := db.Exec(query,
            company.name,
            company.desc,
            company.naviUser,
            fmt.Sprintf("QID-%d", i+1),  // Unique query ID
            baseTime.Add(time.Duration(i*24)*time.Hour), // Staggered start times
            10 + i,     // Different durations
            5 + i,      // Different repetitions
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
            trpl_type_name,
            trpl_name,
            balance_begin,
            balance_end,
            subs_status_name,
            subs_device_name,
            region,
            sms_tj,
            sms_ru,
            sms_eng
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `

    regions := []string{"Dushanbe", "Khujand", "Kulob", "Bokhtar", "Rasht", "GBAO", "Hisor", "Tursunzoda", "Panjakent", "Istaravshan"}
    deviceTypes := []string{"Android", "iOS", "Web", "Desktop", "Tablet", "Smart TV", "Mobile Web", "PWA", "Terminal", "POS"}
    statuses := []string{"Active", "Pending", "Suspended", "Trial", "Expired", "New", "VIP", "Standard", "Premium", "Basic"}

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

    for i, companyID := range companyIDs {
        _, err := db.Exec(query,
            companyID,
            fmt.Sprintf("Type %d", i+1),
            fmt.Sprintf("Triple Name %d", i+1),
            float64(50+i*10),    // balance_begin
            float64(100+i*20),   // balance_end
            statuses[i],
            deviceTypes[i],
            regions[i],
            fmt.Sprintf("SMS TJ Template %d", i+1),
            fmt.Sprintf("SMS RU Template %d", i+1),
            fmt.Sprintf("SMS ENG Template %d", i+1),
        )
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
        efficiency := 0.5 + float64(i)*0.05  // Efficiency from 0.5 to 0.95
        subAmount := 1000 + i*100            // Sub amount from 1000 to 1900
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