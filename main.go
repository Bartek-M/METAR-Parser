package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	API      string
	Interval int
	Stations []string
	Exclude  []string
}

type Weather struct {
	station   string
	time      string
	windDir   string
	windSpeed string
	qnh       string
	metar     string
}

func openConfig() (*Config, error) {
	file, err := os.Open("./config.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	var data Config

	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
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

func parseMetar(val string) (Weather, error) {
	metar := strings.Split(val, " ")

	reQNH := regexp.MustCompile(`Q(\d{4})`)
	matchQNH := reQNH.FindString(val)

	reWind := regexp.MustCompile(`(\d{3})(\d{2})(G(\d{2}))?KT`)
	matchWind := reWind.FindStringSubmatch(val)

	return Weather{
		station:   metar[0],
		time:      metar[1],
		windDir:   matchWind[1],
		windSpeed: matchWind[2],
		qnh:       matchQNH,
		metar:     val,
	}, nil
}

func filterData(data map[string]Weather, exclude []string) {
	for _, val := range exclude {
		delete(data, val)
	}
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("[ERROR] %s: %v\n", message, err)
		os.Exit(1)
	}
}

func main() {
	fmt.Printf("[METAR]\n\n")

	config, err := openConfig()
	handleError(err, "Failed to load configuration")

	fetched, err := fetchData(config.Stations[0], config.API)
	handleError(err, "Failed to fetch API data")

	data := make(map[string]Weather)
	for _, val := range fetched {
		parsed, err := parseMetar(val)
		handleError(err, "Failed to parse METAR")

		data[parsed.station] = parsed
	}

	filterData(data, config.Exclude)
	fmt.Printf("Data: %v", data)
}
