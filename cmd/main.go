package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"METAR-Parser/internal/config"
	"METAR-Parser/internal/metar"
	"METAR-Parser/internal/runways"
	"METAR-Parser/internal/types"
)

func fetchData(station string, url string) ([]string, error) {
	resp, err := http.Get(url + station)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(body), "\n"), nil
}

func filterData(data map[string]types.Weather, airports map[string]types.Airport) {
	for station := range data {
		if _, exists := airports[station]; exists {
			continue
		}
		delete(data, station)
	}
}

func outputData(data map[string]types.Weather, minimums types.Minimums) {
	stations := make([]string, 0, len(data))
	for station := range data {
		stations = append(stations, station)
	}

	sort.Strings(stations)
	for _, station := range stations {
		weather := data[station]

		qnhInfo := ""
		if weather.Qnh != weather.LastQnh && weather.LastQnh != "" {
			qnhInfo = fmt.Sprintf("| %s -> %s", weather.LastQnh, weather.Qnh)
		}

		category := ""
		if weather.Category >= 0 && weather.Category < len(minimums.Category) {
			category = minimums.Category[weather.Category]
		}

		fmt.Printf("%5s | %2s/%2s | %s %s\n", category, weather.DepRwy, weather.ArrRwy, weather.Metar, qnhInfo)
	}
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("[ERROR] %s: %v\n", message, err)
		bufio.NewScanner(os.Stdin).Scan()
		os.Exit(1)
	}
}

func main() {
	conf, err := config.OpenConfig()
	handleError(err, "Failed to load configuration")
	
	data := make(map[string]types.Weather)
	var lastMetars map[string]string
	
	for {
		now := time.Now().UTC()
		fmt.Print("\x1B[2J\x1B[1;1H") // clear terminal
		fmt.Printf("METAR Parser - %s\n\n", fmt.Sprintf("%02d%02d%d0Z", now.Day(), now.Hour(), now.Minute()/10))

		metars := make(map[string]string)
		for _, station := range conf.Stations {
			fetched, err := fetchData(station, conf.API)
			handleError(err, "Failed to fetch API data")
			
			for _, weather := range fetched {
				icao := strings.Split(weather, " ")[0]
				metars[icao] = weather
			}
		}

		for station, weather := range metars {
			if weather == lastMetars[station] {
				continue
			}

			parsed, err := metar.ParseMetar(weather, conf.Minimums)
			handleError(err, "Failed to parse METAR")

			if _, exists := data[station]; exists {
				parsed.LastQnh = data[station].Qnh
			}
			
			runways.AssignRunways(parsed, conf.WindLimits, conf.Airports)
			data[station] = *parsed
		}

		if conf.ExcludeNoConfig {
			filterData(data, conf.Airports)
		}
		outputData(data, conf.Minimums)
		lastMetars = metars

		if conf.Interval == -1 {
			break
		} 

		time.Sleep(time.Duration(conf.Interval * int(time.Second)))
	}

	bufio.NewScanner(os.Stdin).Scan()
}
