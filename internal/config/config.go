package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Environment string `json:"environment"`
	DocumentDatabaseUrl string `json:"documentDatabaseUrl"`
}

func InitializeConfig(fileName string) (Config, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return Config{}, err
	}

	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}