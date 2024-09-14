package types

import "time"

type Company struct {
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	StartDate  time.Time `json:"startDate"`
	Duration   int       `json:"duration"`
	RegionID   int       `json:"regionId"`
	ServiceID  int       `json:"serviceID"`
	TRIPLID    int       `json:"trplId"`
	BalanceMax float64   `json:"balaceMin"`
	BalanceMin float64   `json:"balaceMax"`	
}
