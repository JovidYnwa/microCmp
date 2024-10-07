package types

type PaginatedResponse struct {
	TotalCount  int         `json:"totalCount"`
	TotalPages  int         `json:"totalPages"`
	CurrentPage int         `json:"currentPage"`
	PageSize    int         `json:"pageSize"`
	Data        interface{} `json:"data"`
}

type Company struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	CmpLaunched     int     `json:"cmpLaunched"`
	SubscriberCount int     `json:"subscriberCount"`
	Efficiency      float64 `json:"efficiency"`
}

type CompanyInfo struct {
	CompanyID int `json:"cmp_id"`

	TrplType       int        `json:"trplType"`
	TrplTypeName   string     `json:"trplTypeName"`
	Trpl           BaseFilter `json:"trpl"`
	BalanceBegin   float64    `json:"balanceBegin"`
	BalanceEnd     float64    `json:"balanceEnd"`
	SubsStatusID   int        `json:"subsStatusId"`
	SubsStatusName string     `json:"subsStatusName"`
	SubsDevice     int        `json:"subsDeviceId"`
	SubsDeviceName string     `json:"subsDeviceName"`
	RegionID       int        `json:"regionId"`
	RegionName     string     `json:"regionName"`
	SmsTj          string     `json:"smsTj"`
	SmsRus         string     `json:"smsRus"`
	SmsEng         string     `json:"SmsEng"`
}

type CreateCompanyReq struct {
	Company     Company     `json:"company"`
	CompanyInfo CompanyInfo `json:"companyInfo"`
}
