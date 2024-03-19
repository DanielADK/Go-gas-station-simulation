package configs

import (
	"encoding/json"
	"os"
)

type Config struct {
	Stations  []StationConfig  `json:"stations"`
	Registers []RegisterConfig `json:"registers"`
	Generator GeneratorConfig  `json:"generator"`
}

type StationConfig struct {
	Id        int        `json:"id"`
	Fuel      string     `json:"fuel"`
	TimeRange [2]float64 `json:"time"`
}

type RegisterConfig struct {
	Id        int        `json:"id"`
	TimeRange [2]float64 `json:"time"`
}

type GeneratorConfig struct {
	CountOfCars uint       `json:"count_of_cars"`
	TimeRange   [2]float64 `json:"generation_delay"`
	FuelTypes   []string   `json:"fuel_types"`
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
