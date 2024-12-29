package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Weather struct {
	station   string
	time      string
	windDir   int
	windSpeed int
	qnh       string
	// lastQnh   string ""
	depRwy string
	arrRwy string
	metar  string
}

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

func parseMetar(metar string) (*Weather, error) {
	splitMetar := strings.Split(metar, " ")

	reQNH := regexp.MustCompile(`Q(\d{4})`)
	reWind := regexp.MustCompile(`(\w{3})(\d{2})(G(\d{2}))?KT`)

	qnh := reQNH.FindString(metar)
	wind := reWind.FindStringSubmatch(metar)
	if qnh == "" || wind == nil {
		return nil, fmt.Errorf("METAR incomplete, failed parsing QNH / wind")
	}

	windDir, _ := strconv.Atoi(wind[1])
	windSpeed, _ := strconv.Atoi(wind[2])

	return &Weather{
		station:   splitMetar[0],
		time:      splitMetar[1],
		windDir:   windDir,
		windSpeed: windSpeed,
		qnh:       qnh,
		depRwy:    "--",
		arrRwy:    "--",
		metar:     metar,
	}, nil
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
		fmt.Printf("%s/%s | %s\n", weather.depRwy, weather.arrRwy, weather.metar)
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
		fmt.Print("\033[H\033[2J") // clear terminal
		fmt.Printf("METAR Parser - %s\n\n", fmt.Sprintf("%d%d%d0Z", now.Day(), now.Hour(), now.Minute()/30*3))

		var metars []string
		for _, station := range config.Stations {
			fetched, err := fetchData(station, config.API)
			handleError(err, "Failed to fetch API data")
			metars = append(metars, fetched...)
		}

		for _, val := range metars {
			parsed, err := parseMetar(val)
			handleError(err, "Failed to parse METAR")

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
