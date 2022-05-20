package main

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Influxdb struct {
		Protocol string `yaml:"protocol"`
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Org      string `yaml:"org"`
		Bucket   string `yaml:"bucket"`
		Token    string `yaml:"token"`
	} `yaml:"influxdb"`
	Speedtest struct {
		Server struct {
			Id   int    `yaml:"id"`
			Name string `yaml:"name"`
		}
	} `yaml:"speedtest"`
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
}

func NewConfig(configPath string) (*Config, error) {
	if err := ValidateConfigPath(configPath); err != nil {
		return nil, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	config := &Config{}
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	config = envConfig(config)

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
	newConfig.Influxdb.Protocol = getEnv("INFLUXDB_PROTOCOL", config.Influxdb.Protocol)
	newConfig.Influxdb.Address = getEnv("INFLUXDB_ADDRESS", config.Influxdb.Address)
	newConfig.Influxdb.Port = getEnvInt("INFLUXDB_PORT", config.Influxdb.Port)
	newConfig.Influxdb.Org = getEnv("INFLUXDB_ORG", config.Influxdb.Org)
	newConfig.Influxdb.Bucket = getEnv("INFLUXDB_BUCKET", config.Influxdb.Bucket)
	newConfig.Influxdb.Token = getEnv("INFLUXDB_TOKEN", config.Influxdb.Token)
	newConfig.Speedtest.Server.Id = getEnvInt("SPEEDTEST_SERVER_ID", config.Speedtest.Server.Id)
	newConfig.Speedtest.Server.Name = getEnv("SPEEDTEST_SERVER_NAME", config.Speedtest.Server.Name)
	newConfig.Logging.Level = getEnv("LOGGING_LEVEL", config.Logging.Level)
	return newConfig
}
