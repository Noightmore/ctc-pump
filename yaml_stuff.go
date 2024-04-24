package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Stations  map[string]StationConfig `yaml:"stations"`
	Registers RegisterConfig           `yaml:"registers"`
	Cars      CarConfig                `yaml:"cars"`
}

type StationConfig struct {
	Count        int    `yaml:"count"`
	ServeTimeMin string `yaml:"serve_time_min"`
	ServeTimeMax string `yaml:"serve_time_max"`
}

type RegisterConfig struct {
	Count         int    `yaml:"count"`
	HandleTimeMin string `yaml:"handle_time_min"`
	HandleTimeMax string `yaml:"handle_time_max"`
}

type CarConfig struct {
	Count          int64  `yaml:"count"`
	ArrivalTimeMin string `yaml:"arrival_time_min"`
	ArrivalTimeMax string `yaml:"arrival_time_max"`
}

func LoadConfig(filename string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func loadYamlData() *Config {
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	return config
}
