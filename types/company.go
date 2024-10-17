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
	CmpDesc    string    `json:"desc"`
	NaviUser   string    `json:"naviUser"`
	Repition   int       `json:"repition"`
	StartTime  time.Time `json:"startTime"`
	Duration   int       `json:"durationDay"`
	Repetition float64   `json:"repetion"`
	DWHID      string
}

type TrafficSpent struct {
	Min int `json:"min"`
	Sms int `json:"sms"`
	MB  int `json:"mb"`
}

type BalanceLimit struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type ARPULimit struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type CompanyInfo struct {
	CompanyID int `json:"cmp_id"`

	PhoneType        BaseFilter   `json:"phoneType"` //тип номеров
	Trpl             BaseFilter   `json:"trpl"` //по тарифу
	BalanceLimits    BalanceLimit `json:"balanceLimits"` //по балансу
	SubscriberStatus BaseFilter   `json:"subscriberStatus"` //По статусу абонента
	DeviceType       int          `json:"deviceType"` //По дивайсу ОС
	PackSpent        TrafficSpent `json:"packSpent"` //По использованию мегабайтов
	ARPULimits       ARPULimit    `json:"arpuLimits"` //По арпу
	Region           BaseFilter   `json:"region"` //По активным услугам
	SimDate          time.Time    `json:"start"` //По новому подключению симкарты
	Service          BaseFilter   `json:"service"` //По активным услугам
	WheelUsage       bool         `json:"usingWheel"` //По использованию колеса подарков
}

type TextType struct {
	Ru  string `json:"ru"`
	Tj  string `json:"tj"`
	Eng string `json:"eng"`
}

type CompanySmsBefore struct {
	SmsText    TextType   `json:"smsText"`
	SmsDay     int        `json:"remiderDay"`
	SmsDayText BaseFilter `json:"remiderText"`
}

type CompanySmsAfter struct {
	Action BaseFilter `json:"action"`
	Sms    TextType   `json:"smsRemider"`
	Prize  BaseFilter `json:"prize"`
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
