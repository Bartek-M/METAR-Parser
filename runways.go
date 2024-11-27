package main

import (
	"math"
)

func assignRunways(weather *Weather, windLimit int, runwayConfig map[string][]Runways) {
	station := weather.station
	if _, exists := runwayConfig[station]; !exists {
		return
	}

	var depRwy string
	var arrRwy string

	for _, rwy := range runwayConfig[station] {
		if weather.windSpeed > windLimit &&
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
