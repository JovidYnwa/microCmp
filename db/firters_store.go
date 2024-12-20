package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JovidYnwa/microCmp/types"
	go_ora "github.com/sijms/go-ora/v2"
)

type CompanyFilterStore interface {
	GetTrpls(ctx context.Context) ([]*types.BaseFilter, error)
	GetRegions() ([]*types.BaseFilter, error)
	GetSubsStatuses() ([]*types.BaseFilter, error)
	GetServs(ctx context.Context) ([]*types.BaseFilter, error)
	GetSimTypes() ([]*types.BaseFilter, error)
}

type DwhFilterStore struct {
	db *sql.DB
}

func NewOracleMainScreenStore(db *sql.DB) *DwhFilterStore {
	return &DwhFilterStore{
		db: db,
	}
}

func (s *DwhFilterStore) GetTrpls(ctx context.Context) ([]*types.BaseFilter, error) {
	query := `
		select 
			c.trpl_id, 
			c.trpl_name 
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

func (s *DwhFilterStore) GetServs(ctx context.Context) ([]*types.BaseFilter, error) {
	var (
		inputIDVal int = 773
		cursor     go_ora.RefCursor
		resultCode int
		resultText string
		resutPos   string
		cmdText    string
	)

	// Acquire a fixed connection from the pool
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed GetServs to acquire fixed connection: %w", err)
	}
	defer conn.Close()

	cmdText = `BEGIN 
        odsadmin.get_actual_services_by_trpl(
            i_trpl_id => :i_trpl_id,
            o_cursor_data => :o_cursor_data,
            o_result => :o_result,
            o_err_msg => :o_err_msg,
            o_error_position => :o_error_position
        ); 
        commit;
    END;`

	// Use ExecContext instead of Exec
	_, err = conn.ExecContext(ctx, cmdText,
		inputIDVal,
		sql.Out{Dest: &cursor},
		sql.Out{Dest: &resultCode},
		sql.Out{Dest: &resultText},
		sql.Out{Dest: &resutPos},
	)

	if err != nil {
		return nil, fmt.Errorf("executing stored procedure Servs err: %w", err)
	}

	defer cursor.Close()

	// Use QueryContext if available in go-ora, otherwise use Query
	rows, err := cursor.Query()
	if err != nil {
		return nil, fmt.Errorf("retrieving data from Servs cursor err: %w", err)
	}
	defer rows.Close()

	resp := make([]*types.BaseFilter, 0)
	for rows.Next_() {
		// Check context cancellation during iteration
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			var serv types.BaseFilter
			err = rows.Scan(&serv.ID, &serv.Name)
			if err != nil {
				return nil, fmt.Errorf("scanning Servs row err: %w", err)
			}
			resp = append(resp, &serv)
		}
	}

	return resp, nil
}

func (s *DwhFilterStore) GetSimTypes() ([]*types.BaseFilter, error) {
	query := `
	select 
       o.pht_id,
       o.phone_type
	from odsadmin.phone_type o 
	where pht_id in (4,5,6,10,15,16,17,18)`

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
