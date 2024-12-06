package db

import (
	"database/sql"
	"fmt"

	"github.com/JovidYnwa/microCmp/types"
)

type WorkerMethod interface {
	GetActiveCompanies() ([]*types.ActiveCmp, error)
	GetActiveCompanyItarations() ([]*types.ActiveCmpIteration, error)
	InsertCmpStatistic(stat types.CmpStatistic) (*types.CmpStatistic, error)
	UpdateIterationStatistic(cmpId int, cmpIter *types.CmpStatistic) error
}

type WorkerStore struct {
	db *sql.DB
}

func NewWorkerStore(db *sql.DB) WorkerMethod {
	return &WorkerStore{
		db: db,
	}
}

func (s *WorkerStore) GetActiveCompanies() ([]*types.ActiveCmp, error) {
	query := `
		select 
			c.cmp_billing_id,
			c.sms_data
		from company c
		where c.end_date > current_date;
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := []*types.ActiveCmp{}
	for rows.Next() {
		company := new(types.ActiveCmp)
		err := rows.Scan(
			&company.BillingID,
			&company.SmsText,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}
	return companies, nil
}

func (s *WorkerStore) GetActiveCompanyItarations() ([]*types.ActiveCmpIteration, error) {
	query := `
		SELECT 
			c.id,
			c.cmp_billing_id, 
			cr.start_date
		FROM company c
		JOIN company_repetion cr 
		ON c.id = cr.company_id
		WHERE c.end_date + INTERVAL '3 days' > NOW(); --cmp should work 3 day even after cmp completeion
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	companies := []*types.ActiveCmpIteration{}
	for rows.Next() {
		company := new(types.ActiveCmpIteration)
		err := rows.Scan(
			&company.ID,
			&company.BillingID,
			&company.ItarationDay,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}
	return companies, nil
}

func (s *WorkerStore) InsertCmpStatistic(stat types.CmpStatistic) (*types.CmpStatistic, error) {
	var totalCount int

	selectQuery := `
        SELECT COUNT(cr.id)
        FROM company_repetion cr
        WHERE cr.company_id = $1
        AND DATE_TRUNC('day', cr.start_date::timestamp) = DATE_TRUNC('day', $2::timestamp)
    `
	err := s.db.QueryRow(selectQuery, stat.BillingID, stat.StartDate).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("checking existing company statistic: %w", err)
	}

	if totalCount > 0 {
		return nil, nil
	}

	insertQuery := `
        INSERT INTO company_repetion
        (company_id, efficiency, sub_amount, start_date)
        VALUES
        ($1, $2, $3, $4)
        RETURNING company_id, efficiency, sub_amount, start_date
    `

	returnedStat := new(types.CmpStatistic)
	err = s.db.QueryRow(insertQuery,
		stat.BillingID,
		stat.Efficiency,
		stat.SubscriberAmount,
		stat.StartDate,
	).Scan(
		&returnedStat.BillingID,
		&returnedStat.Efficiency,
		&returnedStat.SubscriberAmount,
		&returnedStat.StartDate,
	)

	if err != nil {
		return nil, fmt.Errorf("inserting company statistic: %w", err)
	}

	return returnedStat, nil
}

func (s *WorkerStore) UpdateIterationStatistic(cmpId int, cmpIter *types.CmpStatistic) error {
    if cmpIter == nil {
        return fmt.Errorf("nil statistics provided")
    }

    query := `
        UPDATE company_repetion 
        SET 
            efficiency = $1,
            sub_amount = $2
        WHERE company_id = $3
        AND TO_CHAR(start_date AT TIME ZONE 'UTC', 'DD.MM.YYYY') = TO_CHAR($4 AT TIME ZONE 'UTC', 'DD.MM.YYYY')
        RETURNING company_id`

    var updatedID int
    err := s.db.QueryRow(query,
        cmpIter.Efficiency,
        cmpIter.SubscriberAmount,
        cmpId,  
        cmpIter.StartDate,
    ).Scan(&updatedID)

    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("no matching row found for company ID %v on date %v", cmpId, 
                cmpIter.StartDate.Format("2006-01-02"))
        }
        return fmt.Errorf("update failed for company ID %v: %w", cmpIter.BillingID, err)
    }

    return nil
}