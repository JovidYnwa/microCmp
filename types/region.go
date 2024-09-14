package types

type Region struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

var Regions = []Region{
	{
		Name: "Dushanbe",
		ID:   1,
	},
	{
		Name: "GBAO",
		ID:   2,
	},
	{
		Name: "Katlon",
		ID:   3,
	},
}
