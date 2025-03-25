package types

type Minimums struct {
	Category   [4]string `json:"category"`
	Visibility [4]int    `json:"visibility"`
	Ceiling    [4]int    `json:"ceiling"`
}

type Weather struct {
	Station   string
	Time      string
	WindDir   int
	WindSpeed int
	Qnh       string
	LastQnh   string
	Vis       int
	Clouds    []int
	DepRwy    string
	ArrRwy    string
	Category  int
	Metar     string
}
