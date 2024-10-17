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
	GetCompanies(page, pageSize int) (*types.PaginatedResponse, error)
	GetCompanyByID(comID int) ([]*types.CompanyDetailResp, error)

	SetCompany(c types.Company) (*int, error)
	SetCompanyInfo(info types.CompanyInfo, sms types.SmsBefore, action types.CompanyAction) error
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

func (s *PgCompanyStore) GetCompanies(page, pageSize int) (*types.PaginatedResponse, error) {
	// Count total number of companies
	var totalCount int
	err := s.db.QueryRow("select count(c.id) from company c").Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Calculate pagination values
	totalPages := (totalCount + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	// Query for paginated results
	query := `
		select 
			c.id,
			c.cmp_name, 
			count(cr.company_id) as cmp_amount, 
			sum(cr.sub_amount) as subs_amount, 
			avg(cr.efficiency)*100 as effiency
		from company_repetion cr
		join company c 
		on cr.company_id = c.id
		group by c.id
		order by 4 desc
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
			&cmp.CmpLunched,
			&cmp.SubsAmont,
			&cmp.Efficiency,
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

func (s *PgCompanyStore) SetCompany(c types.Company) (*int, error) {
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
		c.DWHID,
		c.StartTime,
		c.Duration,
		c.Repetition,
	).Scan(&compId)

	if err != nil {
		return nil, fmt.Errorf("error inserting company: %v", err)
	}
	return &compId, nil
}

func (s *PgCompanyStore) SetCompanyInfo(info types.CompanyInfo, sms types.SmsBefore, action types.CompanyAction) error {
	query := `
        INSERT INTO company_info (
            company_id,
            cmp_filter,
			sms_data,
			action_data
        ) VALUES ($1, $2::jsonb, $3::jsonb, $4::jsonb)` // Note the ::jsonb type cast

	// Create a map for the filter data
	filterData := map[string]interface{}{
		"phoneType":        info.PhoneType,
		"trpl":             info.Trpl,
		"balanceLimits":    info.BalanceLimits,
		"subscriberStatus": info.SubscriberStatus,
		"deviceType":       info.DeviceType,
		"packSpent":        info.PackSpent,
		"arpuLimits":       info.ARPULimits,
		"region":           info.Region,
		"start":            info.SimDate,
		"service":          info.Service,
		"usingWheel":       info.WheelUsage,
	}

	// Marshal the map to JSON
	filterJsonData, err := json.Marshal(filterData)
	if err != nil {
		return fmt.Errorf("error marshaling filter data: %v", err)
	}

	// Create a map for the filter data
	sendSmsData := map[string]interface{}{
		"smsText":      sms.SmsText,
		"smsDay":       sms.SmsDay,
		"smsTextRemid": sms.SmsTextRemid,
	}

	// Marshal the map to JSON
	sendSmsJsonData, err := json.Marshal(sendSmsData)
	if err != nil {
		return fmt.Errorf("error marshaling sendsms data: %v", err)
	}

	actionSmsData := map[string]interface{}{
		"action":  action.Action,
		"smsText": action.Sms,
		"prize":   action.Prize,
	}

	// Marshal the map to JSON
	actionSmsJsonData, err := json.Marshal(actionSmsData)
	if err != nil {
		return fmt.Errorf("error marshaling action data: %v", err)
	}

	// Use sql.RawBytes to pass JSON data
	_, err = s.db.Exec(
		query,
		info.CompanyID,
		json.RawMessage(filterJsonData), // Convert to RawMessage
		json.RawMessage(sendSmsJsonData),
		json.RawMessage(actionSmsJsonData),
	)

	if err != nil {
		return fmt.Errorf("error inserting company info: %v", err)
	}
	return nil
}

func (s *PgCompanyStore) GetCompanyByID(comID int) ([]*types.CompanyDetailResp, error) {
	query := `
	select 
		(c.efficiency)*100 as eff, 
		c.sub_amount, 
		c.start_date, 
		c.end_date
	from company_repetion c where c.company_id = $1`

	rows, err := s.db.Query(query, comID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	companies := []*types.CompanyDetailResp{}
	for rows.Next() {
		cmp := new(types.CompanyDetailResp)
		err := rows.Scan(
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

	return companies, nil
}
