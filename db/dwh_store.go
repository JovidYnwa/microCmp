package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
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
		SELECT 
			COUNT(*) AS total_participants,
			CASE 
				WHEN COUNT(*) = 0 THEN 0
				ELSE ROUND(
					COUNT(CASE WHEN c.action_committed = 1 THEN 1 END) * 100.0 / COUNT(*), 2
				)
			END AS efficiency_percentage
		FROM cms_user.CAMPAIGN_DETAILS c
		WHERE c.campaign_id = :1
		AND trunc(c.insert_date) = trunc(to_date(:2, 'YYYY-MM-DD'))
		AND c.is_participate = 1
	`
	dateStr := date.Format("2006-01-02")

	cmp := &types.CmpStatistic{
		ID:        cmpId,
		StartDate: date,
	}

	err := s.db.QueryRow(query, cmpId, dateStr).Scan(&cmp.SubscriberAmount, &cmp.Efficiency)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query error: %w", err)
	}
	return cmp, nil
}

func (s *DwhWorkerStore) GetDWHCompanyID(ctx context.Context, params *types.CreateCompanyReq) (*float64, error) {
	var (
		billingId  float64
		resutCode  int
		wheelUsed  string = "0"
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
		wheelUsed = "1"
	}

	// Format dates
	startDate := params.CompanyInfo.SimDate.Format("02-Jan-2006")
	campaignStartDate := params.StartDate.Format("02-Jan-2006")
	campaignEndDate := params.EndDate.Format("02-Jan-2006")

	resultText = strings.Repeat(" ", 255)

	_, err := s.db.ExecContext(ctx, cmdText,
		sql.Named("i_phoneType", idList(params.CompanyInfo.PhoneType)),
		sql.Named("i_trpl", idList(params.CompanyInfo.Trpl)),
		sql.Named("i_balanceStart", params.CompanyInfo.BalanceLimits.Start),
		sql.Named("i_balanceEnd", params.CompanyInfo.BalanceLimits.End),
		sql.Named("i_subscriberStatus", idList(params.CompanyInfo.SubscriberStatus)),
		sql.Named("i_deviceType", idList(params.CompanyInfo.DeviceType)),
		sql.Named("i_packSpentMin", params.CompanyInfo.PackSpent.Min),
		sql.Named("i_packSpentSms", params.CompanyInfo.PackSpent.Sms),
		sql.Named("i_packSpentMb", params.CompanyInfo.PackSpent.MB),
		sql.Named("i_arpuStart", params.CompanyInfo.ARPULimits.Start),
		sql.Named("i_arpuEnd", params.CompanyInfo.ARPULimits.End),
		sql.Named("i_region", idList(params.CompanyInfo.Region)),
		sql.Named("i_startDate", startDate),
		sql.Named("i_act_services", idList(params.CompanyInfo.Service)),
		sql.Named("i_noact_services", nil), // Empty string for no inactive services
		sql.Named("i_wheel_use", wheelUsed),
		sql.Named("i_campaign_start_dt", campaignStartDate),
		sql.Named("i_campaign_end_dt", campaignEndDate),
		sql.Named("o_campaign_id", sql.Out{Dest: &billingId}),
		sql.Named("o_result_code", sql.Out{Dest: &resutCode}),
		sql.Named("o_result_text", sql.Out{Dest: &resultText}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute stored procedure: %w", err)
	}

	if resutCode != 0 {
		return nil, fmt.Errorf("procedure returned error code 1 %d: %s", resutCode, resultText)
	}

	return &billingId, nil
}

// Helper function to safely convert slice to comma-separated string
func idList(l []types.BaseFilter) string {
	if len(l) == 0 {
		return ""
	}
	var ids []string
	for _, v := range l {
		ids = append(ids, strconv.Itoa(v.ID))
	}
	return strings.Join(ids, ",")
}
