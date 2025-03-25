package config

import (
	"encoding/json"
	"os"

	"METAR-Parser/internal/types"
)


func OpenConfig() (*types.Config, error) {
	file, err := os.Open("./config/config.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	var data types.Config

	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	if data.Interval < 20 && data.Interval != -1 {
		data.Interval = 20
	}

	return &data, nil
}
