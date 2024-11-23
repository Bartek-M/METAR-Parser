package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	API string
	Stations []string
}

func openConfig() *Config  {
	file, err := os.Open("./config.json")
	if err != nil {
		fmt.Printf("Error opening config file: %v\n", err)
	}

	defer file.Close()
	
	decoder := json.NewDecoder(file)
	var data Config

	err = decoder.Decode(&data)
	if err != nil {
		fmt.Printf("Error decoding JSON: %v", err)
	}

	return &data
}

func main() {
	fmt.Printf("[METAR]\n\n")

	config := openConfig()
	fmt.Printf("CONFIG\nAPI route: %v\nStations: %v", config.API, config.Stations)
}
