package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	API       string
	Interval  int
	Stations  []string
	WindLimit int
	Airports  map[string]Airport
}

type Airport struct {
	Runways map[string]int
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
