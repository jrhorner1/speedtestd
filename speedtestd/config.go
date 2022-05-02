package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Influxdb struct {
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"influxdb"`
	Speedtest struct {
		Server struct {
			Id   int    `yaml:"id"`
			Name string `yaml:"name"`
		}
		Interval string `yaml:"interval"`
	} `yaml:"speedtest"`
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a file", path)
	}
	return nil
}
