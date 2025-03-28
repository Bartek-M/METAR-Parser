package runways

import (
	"math"

	"METAR-Parser/internal/types"
)

func GetRwy(ids []string, runways []types.Runway) []types.Runway {
	if len(ids) == 0 {
		return runways
	}

	var result []types.Runway
	for _, id := range ids {
		for _, rwy := range runways {
			if rwy.Id == id {
				result = append(result, rwy)
			}
		}
	}

	return result
}

func CheckRwy(weather *types.Weather, rwy types.Runway, windLimit [2]int) bool {
	windDir := weather.WindDir
	if windDir == 0 {
		windDir = rwy.Hdg
	}

	valid := true
	currentDirection := math.Abs(float64((windDir-rwy.Hdg+180)%360-180)) < 90

	if weather.WindSpeed >= windLimit[0] && !currentDirection {
		valid = false
	}
	if weather.Category > 1 { 
		if !rwy.ILS {
			valid = false
		} else {
			valid = true
		}
	} 
	if weather.WindSpeed >= windLimit[1] && currentDirection {
		valid = true 
	}

	return valid
}

func SelectRwy(weather *types.Weather, runways []types.Runway, windLimit [2]int) string {
	for _, rwy := range runways {
		if !CheckRwy(weather, rwy, windLimit) {
			continue
		}
		
		return rwy.Id
	}

	return "--"
}

func AssignRunways(weather *types.Weather, windLimit [2]int, airports map[string]types.Airport) {
	station := weather.Station
	airport, exists := airports[station]
	if !exists {
		return
	}

	weather.DepRwy = SelectRwy(weather, GetRwy(airport.Preference.Dep, airport.Runways), windLimit)
	weather.ArrRwy = SelectRwy(weather, GetRwy(airport.Preference.Arr, airport.Runways), windLimit)
}
