package types

type Company struct {
	CmpID           int     `json:"id"`
	Name            string  `json:"name"`
	CmpLunced       int     `json:"cmpLunch"`
	SubscriberCount int     `json:"count"`
	Efficincy       float64 `json:"efficiency"`
}
