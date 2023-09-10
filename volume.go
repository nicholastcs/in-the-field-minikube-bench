package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Workers    string `yaml:"workers"`
	Iterations string `yaml:"iterations"`
}

func getConf(path string) (*Config, error) {
	read, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(read, &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
