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
	WindLimit       int
	Airports        map[string]Airport
}

type Airport struct {
	Runways []struct {
		Id string
		Hdg int
	}
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
