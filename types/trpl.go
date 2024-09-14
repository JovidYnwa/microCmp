package types

type TRPL struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

var RatePlans = []Region{
	{
		Name: "Salom 50",
		ID:   1,
	},
	{
		Name: "Salom 80",
		ID:   2,
	},
	{
		Name: "Salom 100",
		ID:   3,
	},
}
