package configs

import (
	"encoding/json"
	"os"
)

type Config struct {
	Stations  []StationConfig  `json:"stations"`
	Registers []RegisterConfig `json:"registers"`
}

type StationConfig struct {
	Fuel      string     `json:"fuel"`
	TimeRange [2]float64 `json:"time"`
}

type RegisterConfig struct {
	Id        int        `json:"id"`
	TimeRange [2]float64 `json:"time"`
}

func LoadConfigFromFile(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
