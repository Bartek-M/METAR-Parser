package types

type Airport struct {
	Runways    []Runway
	Preference struct {
		Dep []string
		Arr []string
	}
	LVP [2]int
}

type Runway struct {
	Id  string
	Hdg int
	ILS bool
}