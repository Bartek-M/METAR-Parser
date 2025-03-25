package types

type Minimums struct {
	Category   [4]string
	Visibility [4]int
	Ceiling    [4]int
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
