package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/JovidYnwa/microCmp/types"
	_ "github.com/lib/pq"
)

type CompanyStore interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(*types.Account) error
	GetAccounts() ([]*types.Account, error)
	GetAccountByID(int) (*types.Account, error)

	GetCompanyType(page, pageSize int) (*types.PaginatedResponse, error)
	GetCompany(page, pageSize int) (*types.PaginatedResponse, error)
	GetCompanies(page, pageSize int, cmpType string) (*types.PaginatedResponse, error)

	GetCompanyByID(cmpID int) ([]*types.CompanyDetailResp, error)
	SetCompanyType(c types.Company) (*int, error)
	SetCompany(cmp *types.CreateCompanyReq) error
	//UpdateCompanyIteration(comID int) error
}

type PgCompanyStore struct {
	db *sql.DB
}

func NewPgCompanyStore(db *sql.DB) *PgCompanyStore {
	return &PgCompanyStore{
		db: db,
	}
}

func (s *PgCompanyStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PgCompanyStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		number SERIAL,
		balance SERIAL,
		created_at TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PgCompanyStore) CreateAccount(acc *types.Account) error {
	query := `INSERT INTO account
	(first_name, last_name, number, balance, created_at)
	VALUES($1, $2, $3, $4, $5)`

	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (s *PgCompanyStore) UpdateAccount(*types.Account) error {
	return nil
}

func (s *PgCompanyStore) DeleteAccount(id int) error {
	query := `delete from account a where a.id=$1`
	_, err := s.db.Query(query, id)
	return err
}

func (s *PgCompanyStore) GetAccountByID(id int) (*types.Account, error) {
	query := `select * from account a where a.id=$1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PgCompanyStore) GetAccounts() ([]*types.Account, error) {
	query := `select * from account`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	account := new(types.Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)

	if err != nil {
		return nil, err
	}
	return account, err
}

func (s *PgCompanyStore) GetCompanyType(page, pageSize int) (*types.PaginatedResponse, error) {
	var totalCount int
	countQuery := `
        SELECT COUNT(DISTINCT ct.id) 
        FROM company_type ct
        JOIN company c ON c.company_type_id = ct.id
        JOIN company_repetion cr ON cr.company_id = c.id`

	err := s.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("count query error: %w", err)
	}

	totalPages := (totalCount + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	query := `
        SELECT 
            ct.id,
            ct.cmp_name,
            COUNT(cr.id) AS repetition_count,
            SUM(cr.sub_amount) AS total_sub_amount,
            ROUND(AVG(cr.efficiency)::NUMERIC * 100.0, 2) AS average_efficiency_percentage
        FROM 
            company_repetion cr
            JOIN company c ON cr.company_id = c.id
            JOIN company_type ct ON c.company_type_id = ct.id
        GROUP BY 
            ct.id, ct.cmp_name
        ORDER BY ct.id
        LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("main query error: %w", err)
	}
	defer rows.Close()

	companies := make([]*types.CompanyTypeResp, 0, pageSize)
	for rows.Next() {
		cmp := new(types.CompanyTypeResp)
		err := rows.Scan(
			&cmp.ID,
			&cmp.Name,
			&cmp.CmpLunched,
			&cmp.SubsAmount,
			&cmp.Efficiency,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		companies = append(companies, cmp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return &types.PaginatedResponse{
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
		Data:        companies,
	}, nil
}

func (s *PgCompanyStore) GetCompany(page, pageSize int) (*types.PaginatedResponse, error) {
	// Count total number of companies
	var totalCount int
	err := s.db.QueryRow(`
	    select count(company_id) from company_repetion group by company_id`).Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	fmt.Print(totalCount)

	// Calculate pagination values
	totalPages := (totalCount + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	// Query for paginated results
	query := `
		SELECT 
			cr.company_id,
			c.cmp_desc ->> 'name' AS name,
			c.cmp_desc ->> 'desc' AS description,  
			ROUND(AVG(cr.efficiency)::NUMERIC * 100.0, 2) AS average_efficiency_percentage,
			SUM(cr.sub_amount) AS total_sub_amount,  
			TO_CHAR((c.cmp_desc ->> 'startTime')::TIMESTAMP, 'DD.MM.YYYY') AS start_date,
			TO_CHAR(
				(c.cmp_desc ->> 'startTime')::TIMESTAMP + 
				(c.cmp_desc ->> 'durationDay')::INTEGER * INTERVAL '1 day', 
				'DD.MM.YYYY'
			) AS end_date  
		FROM 
			company_repetion cr
		JOIN 
			company c ON cr.company_id = c.id 
		GROUP BY 
			cr.company_id, 
			c.cmp_desc ->> 'name', 
			c.cmp_desc ->> 'desc', 
			(c.cmp_desc ->> 'startTime')::TIMESTAMP,  
			(c.cmp_desc ->> 'durationDay')::INTEGER
		LIMIT $1 OFFSET $2`
	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := []*types.CompanyResp{}
	for rows.Next() {
		cmp := new(types.CompanyResp)
		err := rows.Scan(
			&cmp.ID,
			&cmp.Name,
			&cmp.CmpDesc,
			&cmp.Efficiency,
			&cmp.SubsAmount,
			&cmp.StartDate,
			&cmp.EndDate,
		)
		if err != nil {
			return nil, err
		}

		companies = append(companies, cmp)
	}

	return &types.PaginatedResponse{
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
		Data:        companies,
	}, nil
}

func (s *PgCompanyStore) GetCompanies(page, pageSize int, cmpType string) (*types.PaginatedResponse, error) {
	// Count total number of companies
	var totalCount int
	err := s.db.QueryRow(`
		SELECT 
			COUNT(id) 
			FROM company WHERE company_type_id = $1`, cmpType).Scan(&totalCount)
	
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	query := `
		SELECT 
			c.id,
			c.cmp_desc ->> 'name' AS name,
			c.cmp_desc ->> 'desc' AS description,
			ROUND(COALESCE(AVG(cr.efficiency)::NUMERIC * 100.0, 0), 2) AS average_efficiency_percentage,
			COALESCE(SUM(cr.sub_amount), 0) AS total_sub_amount,
			TO_CHAR(COALESCE((c.cmp_desc ->> 'startTime')::TIMESTAMP, NOW()), 'DD.MM.YYYY') AS start_date,
			TO_CHAR(
				COALESCE((c.cmp_desc ->> 'startTime')::TIMESTAMP, NOW()) + 
				COALESCE((c.cmp_desc ->> 'durationDay')::INTEGER, 0) * INTERVAL '1 day',
				'DD.MM.YYYY'
			) AS end_date
		FROM 
			company c
		LEFT JOIN 
			company_repetion cr ON c.id = cr.company_id
		WHERE c.company_type_id = $3
		GROUP BY 
			c.id,                                
			c.cmp_desc ->> 'name',
			c.cmp_desc ->> 'desc',
			(c.cmp_desc ->> 'startTime')::TIMESTAMP,
			(c.cmp_desc ->> 'durationDay')::INTEGER
		LIMIT $1 OFFSET $2
	`
	rows, err := s.db.Query(query, pageSize, offset, cmpType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := []*types.CompanyResp{}
	for rows.Next() {
		cmp := new(types.CompanyResp)
		err := rows.Scan(
			&cmp.ID,
			&cmp.Name,
			&cmp.CmpDesc,
			&cmp.Efficiency,
			&cmp.SubsAmount,
			&cmp.StartDate,
			&cmp.EndDate,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, cmp)
	}

	return &types.PaginatedResponse{
		TotalCount:  totalCount,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
		Data:        companies,
	}, nil
}


func (s *PgCompanyStore) SetCompanyType(c types.Company) (*int, error) {
	var compId int

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
		RETURNING id`

	err := s.db.QueryRow(
		query,
		c.CmpName,
		c.CmpDesc,
		c.NaviUser,
		// c.StartTime,
		// c.Duration,
		// c.Repetition,
	).Scan(&compId)

	if err != nil {
		return nil, fmt.Errorf("error inserting company: %v", err)
	}
	return &compId, nil
}

func (s *PgCompanyStore) SetCompany(cmp *types.CreateCompanyReq) error {
	query := `
        INSERT INTO company (
            company_type_id,
			start_date,
			end_date,
			cmp_billing_id,
			cmp_desc,
            cmp_filter,
			sms_data,
			action_data
        ) VALUES ($1, $2, $3, $4, $5::jsonb, $6::jsonb, $7::jsonb, $8::json)` // Note the ::jsonb type cast

	cmpData := map[string]interface{}{
		"name":     cmp.Company.CmpName,
		"desc":     cmp.Company.CmpDesc,
		"naviUser": cmp.Company.NaviUser, //Take form token letter on
	}

	cmpJsonData, err := json.Marshal(cmpData)
	if err != nil {
		return fmt.Errorf("error marshaling filter data: %v", err)
	}

	filterData := map[string]interface{}{
		"phoneType":        cmp.CompanyInfo.PhoneType,
		"trpl":             cmp.CompanyInfo.Trpl,
		"balanceLimits":    cmp.CompanyInfo.BalanceLimits,
		"subscriberStatus": cmp.CompanyInfo.SubscriberStatus,
		"deviceType":       cmp.CompanyInfo.DeviceType,
		"packSpent":        cmp.CompanyInfo.PackSpent,
		"arpuLimits":       cmp.CompanyInfo.ARPULimits,
		"region":           cmp.CompanyInfo.Region,
		"start":            cmp.CompanyInfo.SimDate,
		"service":          cmp.CompanyInfo.Service,
		"usingWheel":       cmp.CompanyInfo.WheelUsage,
	}

	filterJsonData, err := json.Marshal(filterData)
	if err != nil {
		return fmt.Errorf("error marshaling filter data: %v", err)
	}

	sendSmsData := map[string]interface{}{
		"smsText":      cmp.SendSms.SmsText,
		"smsDay":       cmp.SendSms.SmsDay,
		"smsTextRemid": cmp.SendSms.SmsTextRemid,
	}

	sendSmsJsonData, err := json.Marshal(sendSmsData)
	if err != nil {
		return fmt.Errorf("error marshaling sendsms data: %v", err)
	}

	actionSmsData := map[string]interface{}{
		"action":  cmp.Action.Action,
		"smsText": cmp.Action.Sms,
		"prize":   cmp.Action.Prize,
	}

	actionSmsJsonData, err := json.Marshal(actionSmsData)
	if err != nil {
		return fmt.Errorf("error marshaling action data: %v", err)
	}

	_, err = s.db.Exec(
		query,
		cmp.CompanyType,
		cmp.StartDate.Time,
		cmp.EndDate.Time,
		cmp.CmpBillingID,
		json.RawMessage(cmpJsonData),
		json.RawMessage(filterJsonData),
		json.RawMessage(sendSmsJsonData),
		json.RawMessage(actionSmsJsonData),
	)

	if err != nil {
		return fmt.Errorf("error inserting company info: %v", err)
	}
	return nil
}

func (s *PgCompanyStore) GetCompanyByID(cmpID int) ([]*types.CompanyDetailResp, error) {
	query := `
		select cr.id,
			cr.efficiency,
			cr.sub_amount,
			TO_CHAR(cr.start_date::TIMESTAMP, 'dd.mm.yyyy')
			--TO_CHAR(cr.end_date::TIMESTAMP, 'dd.mm.yyyy')
		from company_repetion cr
		where cr.company_id = $1
	`

	rows, err := s.db.Query(query, cmpID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	companies := []*types.CompanyDetailResp{}
	for rows.Next() {
		cmp := new(types.CompanyDetailResp)
		err := rows.Scan(
			&cmp.ID,
			&cmp.Efficiency,
			&cmp.SubsAmount,
			&cmp.StartDate,
			//&cmp.EndDate,
		)
		if err != nil {
			return nil, err
		}

		companies = append(companies, cmp)
	}

	return companies, nil
}
