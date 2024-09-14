package types

type Service struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

var Services = []Region{
	{
		Name: "10 Gb",
		ID:   1,
	},
	{
		Name: "20 GB",
		ID:   2,
	},
	{
		Name: "30 Gb",
		ID:   3,
	},
}
