package types

type Airport struct {
	Runways    []Runway `json:"runways"`
	Preference struct {
		Dep []string `json:"dep"`
		Arr []string `json:"arr"`
	} `json:"preference"`
	LVP [2]int `json:"lvp"`
}

type Runway struct {
	Id  string `json:"id"`
	Hdg int    `json:"hdg"`
	ILS bool   `json:"ils"`
}
