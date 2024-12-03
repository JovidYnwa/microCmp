package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/JovidYnwa/microCmp/types"
)

type DwhStore interface {
	GetCmpSubscribersNotify(cmpID int) ([]*types.CmpSubscriber, error)
	GetCompanyStatistic(cmpId int, date time.Time) (*types.CmpStatistic, error)
	GetDWHCompanyID(ctx context.Context, params *types.CreateCompanyReq) (*float64, error)
}

type DwhWorkerStore struct {
	db *sql.DB
}

func NewDwhWorkerStore(db *sql.DB) *DwhWorkerStore {
	return &DwhWorkerStore{
		db: db,
	}
}

func (s *DwhWorkerStore) GetCmpSubscribersNotify(cmpID int) ([]*types.CmpSubscriber, error) {
	query := `
    select c.msisdn, c.lang_id
    from cms_user.CAMPAIGN_DETAILS c
    where c.campaign_id = :cmpID
    and c.action_committed = 0
    and c.notified_count < 2
	`
	//dateStr := cmpDate.Format("2006-01-02")

	rows, err := s.db.Query(query, cmpID)
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

func (s *DwhWorkerStore) GetDWHCompanyID(ctx context.Context, params *types.CreateCompanyReq) (*float64, error) {
	var (
		billingId  float64
		resutCode  int
		wheelUsed  int = 0
		resultText string
		cmdText    string
	)

	cmdText = `BEGIN
         cms_modules.get_cms_campaign_id(
			i_phoneType => :i_phoneType,
			i_trpl => :i_trpl,
			i_balanceStart => :i_balanceStart,
			i_balanceEnd => :i_balanceEnd,
			i_subscriberStatus => :i_subscriberStatus,
			i_deviceType => :i_deviceType,
			i_packSpentMin => :i_packSpentMin,
			i_packSpentSms => :i_packSpentSms,
			i_packSpentMb => :i_packSpentMb,
			i_arpuStart => :i_arpuStart,
			i_arpuEnd => :i_arpuEnd,
			i_region => :i_region,
			i_startDate => :i_startDate,
			i_act_services => :i_act_services,
			i_noact_services => :i_noact_services,
			i_wheel_use => :i_wheel_use,
			i_campaign_start_dt => :i_campaign_start_dt,
			i_campaign_end_dt => :i_campaign_end_dt,
			o_campaign_id => :o_campaign_id,
			o_result_code => :o_result_code,
			o_result_text => :o_result_text
        );
    END;`

	if params.CompanyInfo.WheelUsage {
		wheelUsed = 1
	} else {
		wheelUsed = 0
	}

	// Use ExecContext instead of Exec
	_, err := s.db.ExecContext(ctx, cmdText,
		params.CompanyInfo.PhoneType[0].ID,
		params.CompanyInfo.Trpl[0].ID,
		params.CompanyInfo.BalanceLimits.Start,
		params.CompanyInfo.BalanceLimits.End,
		params.CompanyInfo.SubscriberStatus[0].ID,
		params.CompanyInfo.DeviceType[0].ID,
		params.CompanyInfo.PackSpent.Min,
		params.CompanyInfo.PackSpent.Sms,
		params.CompanyInfo.PackSpent.MB,
		params.CompanyInfo.ARPULimits.Start,
		params.CompanyInfo.ARPULimits.End,
		params.CompanyInfo.Region[0].ID,
		params.CompanyInfo.SimDate.String(),
		params.CompanyInfo.Service[0].ID,
		"0",
		string(wheelUsed), //should be dynamic
		params.StartDate.String(),
		params.EndDate.String(),
		sql.Out{Dest: &billingId},
		sql.Out{Dest: &resutCode},
		sql.Out{Dest: &resultText},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get billingID GetDWHCompanyID: %w", err)
	}

	if resutCode != 0 {
		return nil, fmt.Errorf("failed to get billingID GetDWHCompanyID: %w", err)
	}
	fmt.Println(billingId)
	fmt.Println(resutCode)

	return &billingId, nil
}
