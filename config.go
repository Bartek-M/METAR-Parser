package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	API             string
	Interval        int
	Stations        []string
	ExcludeNoConfig bool
	Minimums        Minimums
	WindLimit       [2]int
	Airports        map[string]Airport
}

type Minimums struct {
	Category   [4]string
	Visibility [4]int
	Ceiling    [4]int
}

type Airport struct {
	Runways []Runway
	Preference struct {
		Dep []string
		Arr []string
	}
	LVP [2]int
}

type Runway struct {
	Id  string
	Hdg int
	ILS bool
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
