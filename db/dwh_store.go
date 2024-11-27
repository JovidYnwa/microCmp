package db

import (
	"database/sql"
	"time"

	"github.com/JovidYnwa/microCmp/types"
)

type DwhStore interface {
	GetCmpSubscribersNotify(cmpId int) ([]*types.CmpSubscriber, error)
	GetCompanyStatistic(cmpId int, date time.Time) (*types.CmpStatistic, error)
}

type DwhWorkerStore struct {
	db *sql.DB
}

func NewDwhWorkerStore(db *sql.DB) *DwhWorkerStore {
	return &DwhWorkerStore{
		db: db,
	}
}

func (s *DwhWorkerStore) GetCmpSubscribersNotify(cmpId int, cmpDate time.Time) ([]*types.CmpSubscriber, error) {
	query := `
		select c.msisdn, c.lang_id
		from cms_user.CAMPAIGN_DETAILS c
		where c.campaign_id =:1
		and trunc(c.insert_date) = trunc(to_date(:2, 'YYYY-MM-DD'))
		and c.action_committed = 0
		and c.notified_count < 2;
	`
	dateStr := cmpDate.Format("2006-01-02")

	rows, err := s.db.Query(query, cmpId, dateStr)
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

func (s *DwhWorkerStore) GetCompanyStatistic(cmpId int, date time.Time) (*types.CmpStatistic, error) {
	query := `
        select 
            count(c.campaign_id)
        from cms_user.CAMPAIGN_DETAILS c 
        where c.is_participate = 1
        and c.campaign_id = :1
        and trunc(c.insert_date) = trunc(to_date(:2, 'YYYY-MM-DD'))
    `
	dateStr := date.Format("2006-01-02")

	cmp := &types.CmpStatistic{
		ID:        cmpId,
		StartDate: date,
	}

	err := s.db.QueryRow(query, cmpId, dateStr).Scan(&cmp.SubscriberAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return cmp, nil
}
