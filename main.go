package main

import (
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

func outputData(data map[string]Weather) {
	stations := make([]string, 0, len(data))
	for station := range data {
		stations = append(stations, station)
	}

	sort.Strings(stations)
	for _, station := range stations {
		weather := data[station]
		fmt.Printf("%s | %s/%s | %s\n", weather.category, weather.depRwy, weather.arrRwy, weather.metar)
	}
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("[ERROR] %s: %v\n", message, err)
		os.Exit(1)
	}
}

func main() {
	config, err := openConfig()
	handleError(err, "Failed to load configuration")

	data := make(map[string]Weather)

	for {
		now := time.Now().UTC()
		spacer := ""
		if now.Day() < 10 {
			spacer = "0"
		}

		fmt.Print("\033[H\033[2J") // clear terminal
		fmt.Printf("METAR Parser - %s\n\n", fmt.Sprintf("%s%d%d%d0Z", spacer, now.Day(), now.Hour(), now.Minute()/10))

		var metars []string
		for _, station := range config.Stations {
			fetched, err := fetchData(station, config.API)
			handleError(err, "Failed to fetch API data")
			metars = append(metars, fetched...)
		}

		for _, val := range metars {
			parsed, err := parseMetar(val)
			handleError(err, "Failed to parse METAR")
			fmt.Printf("%v\n", parsed)

			assignRunways(parsed, config.WindLimit, config.Airports)
			data[parsed.station] = *parsed
		}

		if config.ExcludeNoConfig {
			filterData(data, config.Airports)
		}
		outputData(data)

		if config.Interval == -1 {
			break
		}
		time.Sleep(time.Duration(config.Interval * int(time.Second)))
	}
}
