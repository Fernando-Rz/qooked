package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	TestEnvironment     bool   `json:"testEnvironment"`
	DocumentDatabaseUrl string `json:"documentDatabaseUrl"`
	DatabaseName        string `json:"databaseName"`
}

func NewConfig(fileName string) (Config, error) {
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

func (config *Config) Validate() error {
	if config.DocumentDatabaseUrl == "" {
		return fmt.Errorf("the configuration 'documentDatabaseUrl' is required")
	}

	if config.DatabaseName == "" {
		return fmt.Errorf("the configuration 'databaseName' is required")
	}

	return nil
}
