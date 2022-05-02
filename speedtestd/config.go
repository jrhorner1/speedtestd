package main

import (
	"fmt"
	"os"
	"strconv"

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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		intvalue, _ := strconv.Atoi(value)
		return intvalue
	}
	return fallback
}

func envConfig(config *Config) *Config {
	newConfig := &Config{}
	newConfig.Influxdb.Address = getEnv("INFLUXDB_ADDRESS", config.Influxdb.Address)
	newConfig.Influxdb.Port = getEnvInt("INFLUXDB_PORT", config.Influxdb.Port)
	newConfig.Influxdb.Database = getEnv("INFLUXDB_DATABASE", config.Influxdb.Database)
	newConfig.Influxdb.Username = getEnv("INFLUXDB_USERNAME", config.Influxdb.Username)
	newConfig.Influxdb.Password = getEnv("INFLUXDB_PASSWORD", config.Influxdb.Password)
	newConfig.Speedtest.Server.Id = getEnvInt("SPEEDTEST_SERVER_ID", config.Speedtest.Server.Id)
	newConfig.Speedtest.Server.Name = getEnv("SPEEDTEST_SERVER_NAME", config.Speedtest.Server.Name)
	newConfig.Speedtest.Interval = getEnv("SPEEDTEST_INTERVAL", config.Speedtest.Interval)
	newConfig.Logging.Level = getEnv("LOGGING_LEVEL", config.Logging.Level)
	return newConfig
}
