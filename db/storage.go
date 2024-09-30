package db

import (
	"database/sql"
	"fmt"

	"github.com/JovidYnwa/microCmp/types"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(*types.Account) error
	GetAccounts() ([]*types.Account, error)
	GetAccountByID(int) (*types.Account, error)
	GetCompanies(page, pageSize int) (*types.PaginatedResponse, error)
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
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

func (s *PostgresStore) CreateAccount(acc *types.Account) error {
	query := `INSERT INTO account
	(first_name, last_name, number, balance, created_at)
	VALUES($1, $2, $3, $4, $5)`

	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt)

	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)
	return nil
}

func (s *PostgresStore) UpdateAccount(*types.Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `delete from account a where a.id=$1`
	_, err := s.db.Query(query, id)
	return err
}

func (s *PostgresStore) GetAccountByID(id int) (*types.Account, error) {
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

func (s *PostgresStore) GetAccounts() ([]*types.Account, error) {
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

func (s *PostgresStore) GetCompanies(page, pageSize int) (*types.PaginatedResponse, error) {
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
	query := `SELECT *
              FROM company 
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
			&cmp.Name,
			&cmp.CmpLaunched,
			&cmp.SubscriberCount,
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
