package types

type Config struct {
	API             string
	Interval        int
	Stations        []string
	ExcludeNoConfig bool
	Minimums        Minimums
	WindLimits      [2]int
	Airports        map[string]Airport
}