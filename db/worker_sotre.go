package db

import (
	"database/sql"

	"github.com/JovidYnwa/microCmp/types"
)

type WorkerMethod interface {
	GetActiveCompanies() ([]*types.ActiveCmp, error)
	GetCompanySubscribers(cmpId int) ([]*types.CmpSubscriber, error)
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
			&company.ID,
			&company.SmsText,
		)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}
	return companies, nil
}

func (s *WorkerStore) GetCompanySubscribers(cmpId int) ([]*types.CmpSubscriber, error) {
	query := `
		select 
			c.msisdn, 
			c.lang_id
		from cms_user.CAMPAIGN_DETAILS c 
		where c.campaign_id = 7;
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscribers := []*types.CmpSubscriber{}
	for rows.Next() {
		subscriber := new(types.CmpSubscriber)
		err := rows.Scan(
			&subscriber.Msisdn,
			&subscriber.LangID,
		)
		if err != nil {
			return nil, err
		}
		subscribers = append(subscribers, subscriber)
	}
	return subscribers, nil
}
