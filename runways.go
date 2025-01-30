package main

import (
	"math"
)

func assignRunways(weather *Weather, windLimit [2]int, airports map[string]Airport) {
	station := weather.station
	if _, exists := airports[station]; !exists {
		return
	}

	depRwy := "--"
	arrRwy := "--"

	for _, rwy := range airports[station].Runways {
		if weather.windSpeed > windLimit[0] &&
			math.Abs(float64((weather.windDir-rwy.Hdg+180)%360-180)) > 90 {
			continue
		}

		depRwy = rwy.Id
		arrRwy = rwy.Id
		break
	}

	weather.depRwy = depRwy
	weather.arrRwy = arrRwy
}
