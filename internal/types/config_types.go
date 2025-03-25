package types

type Config struct {
	API             string             `json:"api"`
	Interval        int                `json:"interval"`
	Stations        []string           `json:"stations"`
	ExcludeNoConfig bool               `json:"excludeNoConfig"`
	Minimums        Minimums           `json:"minimums"`
	WindLimits      [2]int             `json:"windLimits"`
	Airports        map[string]Airport `json:"airports"`
}
