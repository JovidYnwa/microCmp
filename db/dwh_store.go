package db

import (
	"database/sql"

	"github.com/JovidYnwa/microCmp/types"
)

type DwhStore interface {
	GetCompanySubscribers(cmpId int) ([]*types.CmpSubscriber, error)
}

type DwhWorkerStore struct {
	db *sql.DB
}

func NewDwhWorkerStore(db *sql.DB) *DwhWorkerStore {
	return &DwhWorkerStore{
		db: db,
	}
}

func (s *DwhWorkerStore) GetCompanySubscribers(cmpId int) ([]*types.CmpSubscriber, error) {
	query := `
		select 
			c.msisdn, 
			c.lang_id
		from cms_user.CAMPAIGN_DETAILS c 
		where c.campaign_id =:cmpID
	`

	rows, err := s.db.Query(query, cmpId)
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
