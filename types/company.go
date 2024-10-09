package types

import (
	"time"
)

type PaginatedResponse struct {
	TotalCount  int         `json:"totalCount"`
	TotalPages  int         `json:"totalPages"`
	CurrentPage int         `json:"currentPage"`
	PageSize    int         `json:"pageSize"`
	Data        interface{} `json:"data"`
}

type Company struct {
	ID         int       `json:"id"`
	CmpName    string    `json:"name"`
	NaviUser   string    `json:"naviUser"`
	Repition   int       `json:"repition"`
	StartTime  time.Time `json:"startTime"`
	Duration   int       `json:"durationDay"`
	Repetition float64   `json:"repetion"`
	DWHID      string
}

type CompanyInfo struct {
	CompanyID int `json:"cmp_id"`

	TrplType         int        `json:"trplType"`
	TrplTypeName     string     `json:"trplTypeName"`
	Trpl             BaseFilter `json:"trpl"`
	BalanceBegin     float64    `json:"balanceBegin"`
	BalanceEnd       float64    `json:"balanceEnd"`
	SubscriberStatus BaseFilter `json:"subscriberStatus"`
	SubsDevice       int        `json:"subsDeviceId"`
	SubsDeviceName   string     `json:"subsDeviceName"`
	Region           BaseFilter `json:"region"`
	SmsTj            string     `json:"smsTj"`
	SmsRus           string     `json:"smsRus"`
	SmsEng           string     `json:"SmsEng"`
}

type CreateCompanyReq struct {
	Company     Company     `json:"company"`
	CompanyInfo CompanyInfo `json:"companyInfo"`
}

type CompanyDetail struct {
	Efficiency       float64   `json:"efficiency"`
	SubscriberAmount int64     `json:"subsAmount"`
	StartDate        time.Time `json:"startDate"`
	EndDate          time.Time `json:"endDate"`
}
