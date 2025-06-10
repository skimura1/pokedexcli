package types

type Config struct {
	PageNumber int
}

type Location struct {
	Name string `json: "name"`
	Url string `json: "url"`
}

type LocationArea struct {
	Count int `json: "count"`
	Next string `json: "next"`
	Previous string `json: "previous"`
	Results []Location `json: "results"`
}
