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

func filterData(data map[string]Weather, airports map[string]Airport) {
	for station := range data {
		if _, exists := airports[station]; exists {
			continue
		}
		delete(data, station)
	}
}

func outputData(data map[string]Weather, minimums Minimums) {
	stations := make([]string, 0, len(data))
	for station := range data {
		stations = append(stations, station)
	}

	sort.Strings(stations)
	for _, station := range stations {
		weather := data[station]

		qnhInfo := ""
		if weather.qnh != weather.lastQnh && weather.lastQnh != "" {
			qnhInfo = fmt.Sprintf("%s -> %s | ", weather.qnh, weather.lastQnh)
		}

		category := ""
		if weather.category >= 0 && weather.category < len(minimums.Category) {
			category = minimums.Category[weather.category]
		}

		fmt.Printf("%5s | %2s/%2s | %s%s\n", category, weather.depRwy, weather.arrRwy, qnhInfo, weather.metar)
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
	config, err := openConfig()
	handleError(err, "Failed to load configuration")

	data := make(map[string]Weather)
	var lastMetars map[string]string

	for {
		now := time.Now().UTC()
		fmt.Print("\033[H\033[2J") // clear terminal
		fmt.Printf("METAR Parser - %s\n\n", fmt.Sprintf("%02d%02d%d0Z", now.Day(), now.Hour(), now.Minute()/10))

		metars := make(map[string]string)
		for _, station := range config.Stations {
			fetched, err := fetchData(station, config.API)
			handleError(err, "Failed to fetch API data")
			
			for _, metar := range fetched {
				icao := strings.Split(metar, " ")[0]
				metars[icao] = metar
			}
		}

		for station, metar := range metars {
			if metar == lastMetars[station] {
				continue
			}

			parsed, err := parseMetar(metar, config.Minimums)
			handleError(err, "Failed to parse METAR")

			if _, exists := data[station]; exists {
				parsed.lastQnh = data[station].qnh
			}
			
			assignRunways(parsed, config.WindLimits, config.Airports)
			data[station] = *parsed
		}

		if config.ExcludeNoConfig {
			filterData(data, config.Airports)
		}
		outputData(data, config.Minimums)
		lastMetars = metars

		if config.Interval == -1 {
			break
		} 

		time.Sleep(time.Duration(config.Interval * int(time.Second)))
	}

	bufio.NewScanner(os.Stdin).Scan()
}
