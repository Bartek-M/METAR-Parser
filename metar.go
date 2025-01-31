package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Weather struct {
	station   string
	time      string
	windDir   int
	windSpeed int
	qnh       string
	lastQnh   string
	vis       int
	clouds    []int
	depRwy    string
	arrRwy    string
	category  int
	metar     string
}

func parseWind(metar string) (int, int, error) {
	reWind := regexp.MustCompile(`(\w{3})(\d{2})(G(\d{2}))?KT`)
	wind := reWind.FindStringSubmatch(metar)

	windDir, errDir := strconv.Atoi(wind[1])
	windSpeed, errSpeed := strconv.Atoi(wind[2])

	if wind[1] != "VRB" && errDir != nil {
		return 0, 0, fmt.Errorf("failed parsing wind direction")
	}
	if errSpeed != nil {
		return 0, 0, fmt.Errorf("failed parsing wind speed")
	}

	return windDir, windSpeed, nil
}

func parseVisibility(metar string) (int, error) {
	if strings.Contains(metar, "CAVOK") {
		return 9999, nil
	}

	reVis := regexp.MustCompile(`\b(\d{4})\b`)
	visMatch := reVis.FindStringSubmatch(metar)

	vis, err := strconv.Atoi(visMatch[0])
	return vis, err
}

func parseClouds(metar string) []int {
	if strings.Contains(metar, "CAVOK") {
		return []int{9999}
	}

	reClouds := regexp.MustCompile(`(?:FEW|SCT|BKN|OVC)(\d{3})`)
	cloudMatches := reClouds.FindAllStringSubmatch(metar, -1)

	if len(cloudMatches) < 1 {
		return []int{9999}
	}

	var clouds []int
	for _, match := range cloudMatches {
		if height, err := strconv.Atoi(match[1]); err == nil {
			clouds = append(clouds, height*100)
		}
	}

	return clouds
}

func parseQNH(metar string) string {
	reQNH := regexp.MustCompile(`Q(\d{4})`)
	qnh := reQNH.FindString(metar)

	return qnh
}

func getCategory(visibility int, clouds []int, minimums Minimums) int {
	for i := range 4 {
		if visibility < minimums.Visibility[i] || clouds[0] < minimums.Ceiling[i] {
			continue
		}

		return i
	}

	return -1
}

func parseMetar(metar string, minimums Minimums) (*Weather, error) {
	splitMetar := strings.Split(metar, " ")

	windDir, windSpeed, err := parseWind(metar)
	if err != nil {
		return nil, fmt.Errorf("failed parsing wind")
	}

	clouds := parseClouds(metar)
	vis, err := parseVisibility(metar)
	if err != nil {
		return nil, fmt.Errorf("failed parsing visibility")
	}

	qnh := parseQNH(metar)
	if qnh == "" {
		return nil, fmt.Errorf("failed parsing QNH")
	}

	return &Weather{
		station:   splitMetar[0],
		time:      splitMetar[1],
		windDir:   windDir,
		windSpeed: windSpeed,
		qnh:       qnh,
		lastQnh:   "",
		vis:       vis,
		clouds:    clouds,
		depRwy:    "--",
		arrRwy:    "--",
		category:  getCategory(vis, clouds, minimums),
		metar:     metar,
	}, nil
}
