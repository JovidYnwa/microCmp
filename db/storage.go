package db

import (
	"database/sql"
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
	GetCompanyByID(comID int) (*types.PaginatedResponse, error)

	SetCompany(c types.Company) (*int, error)
	SetCompanyInfo(c types.CompanyInfo) error
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
	err := s.db.QueryRow("SELECT COUNT(*) FROM company").Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	// Calculate pagination values
	totalPages := (totalCount + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	// Query for paginated results
	query := `SELECT c.id, c.cmp_name, c.start_time, c.duration, c.repetition
              FROM company c
              ORDER BY id 
              LIMIT $1 OFFSET $2`
	rows, err := s.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := []*types.Company{}
	for rows.Next() {
		cmp := new(types.Company)
		err := rows.Scan(
			&cmp.ID,
			&cmp.CmpName,
			&cmp.StartTime,
			&cmp.Duration,
			&cmp.Repetition,
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

// Function to insert company info
func (s *PgCompanyStore) SetCompanyInfo(info types.CompanyInfo) error {
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
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := s.db.Exec(
		query,
		info.CompanyID,
		info.TrplTypeName,
		info.Trpl.Name,
		info.BalanceBegin,
		info.BalanceEnd,
		info.SubscriberStatus.Name,
		info.SubsDeviceName,
		info.Region.Name,
		info.SmsTj,
		info.SmsRus,
		info.SmsEng,
	)

	if err != nil {
		return fmt.Errorf("error inserting company info: %v", err)
	}

	return nil
}

func (s *PgCompanyStore) GetCompanyByID(comID int) (*types.PaginatedResponse, error) {
	query := `SELECT c.id, c.cmp_name, c.start_time, c.duration, c.repetition
	FROM company_repetion c
	ORDER BY id 
	LIMIT $1 OFFSET $2`
	
	rows, err := s.db.Query(query)
	if err != nil {
		fmt.Println("gaga")
	}
	defer rows.Close()

	return nil, nil
}