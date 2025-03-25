package metar

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"METAR-Parser/internal/types"
)


func ParseWind(metar string) (int, int, error) {
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

func ParseVisibility(metar string) (int, error) {
	if strings.Contains(metar, "CAVOK") {
		return 9999, nil
	}

	reVis := regexp.MustCompile(`\b(\d{4})\b`)
	visMatch := reVis.FindStringSubmatch(metar)

	vis, err := strconv.Atoi(visMatch[0])
	return vis, err
}

func ParseClouds(metar string) []int {
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

func ParseQNH(metar string) string {
	reQNH := regexp.MustCompile(`Q(\d{4})`)
	qnh := reQNH.FindString(metar)

	return qnh
}

func GetCategory(visibility int, clouds []int, minimums types.Minimums) int {
	for i := range 4 {
		if visibility < minimums.Visibility[i] || clouds[0] < minimums.Ceiling[i] {
			continue
		}

		return i
	}

	return -1
}

func ParseMetar(metar string, minimums types.Minimums) (*types.Weather, error) {
	splitMetar := strings.Split(metar, " ")

	windDir, windSpeed, err := ParseWind(metar)
	if err != nil {
		return nil, fmt.Errorf("failed parsing wind")
	}

	clouds := ParseClouds(metar)
	vis, err := ParseVisibility(metar)
	if err != nil {
		return nil, fmt.Errorf("failed parsing visibility")
	}

	qnh := ParseQNH(metar)
	if qnh == "" {
		return nil, fmt.Errorf("failed parsing QNH")
	}

	return &types.Weather{
		Station:  splitMetar[0],
		Time:     splitMetar[1],
		WindDir:  windDir,
		WindSpeed: windSpeed,
		Qnh:      qnh,
		LastQnh:  "",
		Vis:      vis,
		Clouds:   clouds,
		DepRwy:   "--",
		ArrRwy:   "--",
		Category: GetCategory(vis, clouds, minimums),
		Metar:    metar,
	}, nil
}
