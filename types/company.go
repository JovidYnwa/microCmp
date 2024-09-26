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