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
	CompanyID    int     `json:"cmp_id"`
	TrplType     string  `json:"trplType"`
	TrplID       int     `json:"trplId"`
	BalanceBegin float64 `json:"balanceBegin"`
	BalanceEnd   float64 `json:"balanceEnd"`
	SubsStatus   int     `json:"subsStatus"`
	SubsDevice   int     `json:"subsDevice"`
	Region       int     `json:"region"`
	SmsTj        string  `json:"smsTj"`
	SmsRus       string  `json:"smsRus"`
	SmsEng       string  `json:"SmsEng"`
}

type CreateCompanyReq struct {
	Company     Company     `json:"company"`
	CompanyInfo CompanyInfo `json:"companyInfo"`
}
