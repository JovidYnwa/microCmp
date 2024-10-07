package db

import (
	"database/sql"
	"fmt"

	"github.com/JovidYnwa/microCmp/types"
)

type CompanyFilterStore interface {
	GetTrpls() ([]*types.BaseFilter, error)
	GetRegions() ([]*types.BaseFilter, error)
	GetSubsStatuses() ([]*types.BaseFilter, error)
}

type DwhFilterStore struct {
	db *sql.DB
}

func NewOracleMainScreenStore(db *sql.DB) *DwhFilterStore {
	return &DwhFilterStore{
		db: db,
	}
}

func (s *DwhFilterStore) GetTrpls() ([]*types.BaseFilter, error) {
	query := `select c.trpl_id, c.trpl_name 
				from CMS_USER.current_tp_names c`

	rows, err := s.db.Query(query)
	if err != nil {
		fmt.Println("gaga")
	}
	defer rows.Close()

	trpls := []*types.BaseFilter{}
	for rows.Next() {
		trpl := new(types.BaseFilter)
		err := rows.Scan(
			&trpl.ID,
			&trpl.Name,
		)
		if err != nil {
			return nil, err
		}
		trpls = append(trpls, trpl)
	}
	return trpls, nil
}

func (s *DwhFilterStore) GetRegions() ([]*types.BaseFilter, error) {
	query := `select o.regiongg_id, o.region from ODSADMIN.PQ_AREA o`

	rows, err := s.db.Query(query)
	if err != nil {
		fmt.Println("gaga")
	}
	defer rows.Close()

	regions := []*types.BaseFilter{}
	for rows.Next() {
		region := new(types.BaseFilter)
		err := rows.Scan(
			&region.ID,
			&region.Name,
		)
		if err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}
	return regions, nil
}

func (s *DwhFilterStore) GetSubsStatuses() ([]*types.BaseFilter, error) {
	query := `select o.stat_id, o.status from odsadmin.status o`

	rows, err := s.db.Query(query)
	if err != nil {
		fmt.Println("gaga")
	}
	defer rows.Close()

	regions := []*types.BaseFilter{}
	for rows.Next() {
		region := new(types.BaseFilter)
		err := rows.Scan(
			&region.ID,
			&region.Name,
		)
		if err != nil {
			return nil, err
		}
		regions = append(regions, region)
	}
	return regions, nil
}