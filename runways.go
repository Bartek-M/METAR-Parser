package main

import (
	"math"
)

func getRwy(ids []string, runways []Runway) []Runway {
	if len(ids) == 0 {
		return runways
	}

	var result []Runway
	for _, id := range ids {
		for _, rwy := range runways {
			if rwy.Id == id {
				result = append(result, rwy)
			}
		}
	}

	return result
}

func checkRwy(weather *Weather, rwy Runway, windLimit [2]int) bool {
	currentDirection := math.Abs(float64((weather.windDir-rwy.Hdg+180)%360-180)) > 90

	if weather.windSpeed > windLimit[0] && currentDirection {
		return false
	}
	if weather.category > 1 && !rwy.ILS && !(weather.windSpeed > windLimit[1] && currentDirection) {
		return false
	}

	return true
}

func selectRwy(weather *Weather, runways []Runway, windLimit [2]int) string {
	for _, rwy := range runways {
		if !checkRwy(weather, rwy, windLimit) {
			continue
		}

		return rwy.Id
	}

	return "--"
}

func assignRunways(weather *Weather, windLimit [2]int, airports map[string]Airport) {
	station := weather.station
	airport, exists := airports[station]
	if !exists {
		return
	}

	weather.depRwy = selectRwy(weather, getRwy(airport.Preference.Dep, airport.Runways), windLimit)
	weather.arrRwy = selectRwy(weather, getRwy(airport.Preference.Arr, airport.Runways), windLimit)
}
